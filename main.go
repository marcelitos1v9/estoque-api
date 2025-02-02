package main

import (
    "estoque-api/database"
    "estoque-api/handlers"
    "github.com/gin-gonic/gin"
)

func main() {
    database.Connect()
    handlers.InitializeHandlers()

    r := gin.Default()

    // Rotas de Produtos
    produtos := r.Group("/produtos")
    {
        produtos.GET("", handlers.GetProdutos)
        produtos.GET("/:id", handlers.GetProduto)
        produtos.POST("", handlers.CreateProduto)
        produtos.PUT("/:id", handlers.UpdateProduto)
        produtos.DELETE("/:id", handlers.DeleteProduto)
        
        // Novas rotas
        produtos.GET("/categoria/:categoria", handlers.GetProdutosPorCategoria)
        produtos.GET("/busca", handlers.BuscarProdutos)
        produtos.PATCH("/:id/estoque", handlers.AtualizarEstoque)
        produtos.PATCH("/:id/preco", handlers.AtualizarPreco)
        produtos.GET("/baixo-estoque", handlers.GetProdutosBaixoEstoque)
        produtos.POST("/:id/imagem", handlers.UploadImagemProduto)
    }

    // Rotas de Relatórios
    relatorios := r.Group("/relatorios")
    {
        relatorios.GET("/estoque", handlers.RelatorioEstoque)
        relatorios.GET("/produtos-mais-vendidos", handlers.RelatorioProdutosMaisVendidos)
        relatorios.GET("/valor-total-estoque", handlers.RelatorioValorTotalEstoque)
    }

    r.Run(":8080")
}
