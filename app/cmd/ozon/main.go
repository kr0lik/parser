package main

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"parser/di/ozon"
	"parser/internal/config"
	configOzon "parser/internal/config/ozon"
	"syscall"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ozon",
		Short: "Парсинг озона",
		Long:  `Парсинг категорий, товаров.`,
	}

	collectCategoriesCmd = &cobra.Command{
		Use:   "collectCategories",
		Short: "Сохранение котегорий",
		Long: `Сначала нужно сохранить нужные категории.
Сделано для того, что бы каждый раз не проходить по всем категориям.`,
		Run: func(cmd *cobra.Command, args []string) {
			uc, err := ozon.InitialiseCollectCategories(cmd.Context())
			if err != nil {
				log.Fatal(err)
			}

			uc.Run()
		},
	}

	collectProductsCmd = &cobra.Command{
		Use:   "collectProducts",
		Short: "Сохранение товаров",
		Long: `Сохранение товаров из предварительно сохраненных категорий(collectCategories).
Сохраняется инфо из катлога, без переходан на страницу товара.`,
		Run: func(cmd *cobra.Command, args []string) {
			uc, err := ozon.InitialiseCollectProducts(cmd.Context())
			if err != nil {
				log.Fatal(err)
			}

			uc.Run()
		},
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Println("Cancel program on signal:", sig)
		cancel()
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	config.ReadMain()
	configOzon.ReadFilter()
	configOzon.ReadParser()

	rootCmd.AddCommand(collectCategoriesCmd)
	rootCmd.AddCommand(collectProductsCmd)
}
