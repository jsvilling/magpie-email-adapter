package main

import (
	"context"
)

func main() {
	ctx := context.Background()
	listeners := make([]TransactionAvailableListener, 1)
	listeners[0] = CreateFileRepository()

	gmailSrv := SetupGmailSrv(&ctx, &listeners)

	httpServer := SetupHttpServer(TransactionProvider(gmailSrv))
	httpServer.listenAndServe()
}
