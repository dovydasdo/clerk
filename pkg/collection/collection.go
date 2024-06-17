package collection

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/dovydasdo/clerk/config"
	i "github.com/dovydasdo/clerk/pkg/ingestion"
	m "github.com/dovydasdo/clerk/pkg/management"
	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"
)

type Manager interface {
	Init() error
	Start() error
	Stop() error
}

type Collection struct {
	Managers []Manager
	l        *slog.Logger
}

var CollectionInstance *Collection
var once sync.Once

func GetCollection(ctx context.Context, cfg config.CollectionConfig, l *slog.Logger) *Collection {
	once.Do(func() {
		CollectionInstance = &Collection{
			l: l,
		}

		l.Debug("collection", "message", fmt.Sprintf("starting collection with config: %+v", cfg))

		for _, manager := range cfg.Managers {
			switch manager.Type {
			case "rent":
				processors := make([]i.Processor, 0)
				for _, procCfg := range manager.Processors {
					sources := getSources(ctx, procCfg, l)
					stores := getStores(procCfg, l)

					var sch gocron.Scheduler
					var err error
					if procCfg.StateDump == "daily" || procCfg.StateDump == "weekly" {
						sch, err = gocron.NewScheduler()
						if err != nil {
							l.Error("init", "error", err)
						}
					}

					procOpts := i.GetRPOptions(
						i.RPWithSource(sources...),
						i.RPWithSaver(stores[0]), // only one for now
						i.RPWithState(&i.RentState{}),
						i.RPWithStateInitF(i.InitSate),
						i.RPWithStateF(i.UpdateRentState),
						i.RPWithLogger(l),
						i.RPWithCtx(context.Background()),
						i.RPWithId(procCfg.Id),
						i.RPWithDumpInterval(procCfg.StateDump),
						i.RPWithScheduler(sch),
					)
					rp := i.GetRentProcessor(*procOpts)

					processors = append(processors, rp)
				}

				rmOpts := m.GetRMOptions(
					m.RMWithCtx(context.TODO()),
					m.RMWithProcessor(processors...),
				)
				rm := m.GetRentManager(*rmOpts)

				CollectionInstance.Managers = append(CollectionInstance.Managers, rm)
				l.Info("init", "message", fmt.Sprintf("starting manager of type %+v", manager.Type))

			default:
				l.Warn("init", "message", fmt.Sprintf("manager of type %v not implemented", manager.Type))
			}
		}
	})

	return CollectionInstance
}

func (c Collection) Start() error {
	for _, mngr := range c.Managers {
		err := mngr.Init()
		if err != nil {
			return err
		}
	}

	for _, mngr := range c.Managers {
		err := mngr.Start()
		if err != nil {
			return err
		}
	}

	for {
	}

}

func getSources(ctx context.Context, cfg config.Processor, l *slog.Logger) []i.Ingestor {
	sources := make([]i.Ingestor, 0)
	for _, sourceCfg := range cfg.DataSources {
		source := getSource(ctx, sourceCfg, l)
		sources = append(sources, source)
	}

	return sources

}

func getSource(ctx context.Context, cfg config.DataSource, l *slog.Logger) i.Ingestor {
	var source i.Ingestor
	switch cfg.SourceType {
	case "nats":
		url := viper.GetString("NATS_SERVER")
		if url == "" {
			l.Warn("source init", "message", "failed to get url of NATS server, using default(localhost:4222)")
			url = "localhost:4222"
		}

		iOpts := i.GetNIOptions(
			i.NIWithCtx(ctx),
			i.NIWithName(cfg.Name),
			i.NIWithLogger(l),
			i.NIWithStreamName(cfg.Stream),
			i.NIWithSubject(cfg.Subject),
			i.NIWithServerUrl(url),
		)
		source = i.GetNATSIngestor(iOpts)
		l.Debug("collection", "message", fmt.Sprintf("getting ingestor of type %+v", cfg.Type))
	case "nats_cache":
		url := viper.GetString("NATS_SERVER")
		if url == "" {
			l.Warn("source init", "message", "failed to get url of NATS server, using default(localhost:4222)")
			url = "localhost:4222"
		}
		cOpts := i.GetNKVCacheOptions(
			i.NKVCacheWithUrl(url),
			i.NKVCacheWithCtx(ctx),
			i.NKVCacheWithName(cfg.Name),
			i.NKVCacheWithSubject(cfg.Subject),
		)

		source = i.GetNatsKVSync(cOpts)
		l.Debug("collection", "message", fmt.Sprintf("getting ingestor of type %+v", cfg.Type))

	}
	return source
}

func getStores(cfg config.Processor, l *slog.Logger) []i.Saver {
	savers := make([]i.Saver, 0)
	for _, saverCfg := range cfg.Stores {
		saver := getStore(saverCfg, l)
		if saver == nil {
			l.Warn("collection", "message", fmt.Sprintf("failed to init store with: %+v", saverCfg))
		}
		savers = append(savers, saver)
	}

	return savers

}

func getStore(cfg config.Store, l *slog.Logger) i.Saver {
	var saver i.Saver
	var err error
	switch cfg.Type {
	case "turso":
		url := viper.GetString("TURSO_URL")
		if url == "" {
			l.Warn("saver init", "message", "failed to get url of turso server, from env variables")
			return nil
		}

		token := viper.GetString("TURSO_TOKEN")
		if token == "" {
			l.Warn("saver init", "message", "failed to get token for turso, from env variables")
			return nil
		}

		// TODO: save functions based on data type
		sOpts := i.GetTSOptions(
			i.TSWithUrl(url),
			i.TSWithToken(token),
			i.TSWithLogger(l),
			i.TSWithSaveFunc(i.SaveAd),
			i.TSWithSaveFunc(i.SaveLocation),
			i.TSWithSaveFunc(i.SaveRentState),
		)

		saver, err = i.GetTursoSaver(sOpts)
		if err != nil {
			l.Error("collection", "error", err)
			return nil
		}

		l.Info("collection", "message", fmt.Sprintf("getting saver of type %+v", cfg.Type))

	default:
	}
	return saver
}
