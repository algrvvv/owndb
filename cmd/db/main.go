package main

import (
	"fmt"
	"os"

	"github.com/algrvvv/owndb/internal/config"
	"github.com/algrvvv/owndb/internal/database"
	"github.com/algrvvv/owndb/internal/dsl"
	"github.com/algrvvv/owndb/internal/logger"
	"github.com/algrvvv/owndb/internal/storage/binarizer"
	"github.com/algrvvv/owndb/internal/storage/memstore"
	"github.com/algrvvv/owndb/internal/storage/snapshot"
	"github.com/algrvvv/owndb/internal/utils"
	"github.com/algrvvv/owndb/internal/wal"
)

func main() {
	// TODO: примерный план для мейна
	// 1. инициализируем конфигурацию
	// 2. читаем данные из дампа
	// 3. читаем данные из wal
	// (скорее всего там будет набор команд,
	// поэтому нам нужно уже к этому моменту иметь парсер
	// и уметь разбирать запрос на команды для базы данных)
	// 4. применяем wal к данным, которые мы уже прочитали
	// 5. держим данные в памяти и держим сервер работаютщим

	// TODO: запускаем дампера, который сохраняет данные раз в указанный промежуток

	conf := config.MustLoad()
	log := logger.MustInit(conf.LogFile, conf.DebugMode)

	srcSnapFile, err := os.OpenFile(conf.SnapshotFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer srcSnapFile.Close()

	marshaller := binarizer.NewBinaryMarshaller(log)
	snap := snapshot.NewSnapshotManager(marshaller, srcSnapFile)
	initData, err := snap.Read()
	if err != nil {
		panic(err)
	}

	walFile, err := os.OpenFile(conf.WalFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}

	wal := wal.NewWAL(walFile)
	mem := memstore.NewMemStorage(initData)
	intr := dsl.NewInterpreter(snap, mem, wal)

	err = utils.PrepareData(log, mem, wal, intr)
	if err != nil {
		panic(err)
	}

	server := database.NewDBServer(conf.Port, log, mem, intr, wal)

	fmt.Println("start database server on: ", conf.Port)
	if err := server.Listen(); err != nil {
		panic(err)
	}
}
