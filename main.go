package main

import (
    "estoque-api/database"
    "estoque-api/handlers"
    "estoque-api/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    database.Connect()
    handlers.InitializeHandlers()
    handlers.InitializeAuthHandlers()

    r := gin.Default()

    // Rotas públicas
    r.POST("/register", handlers.Register)
    r.POST("/login", handlers.Login)

    // Grupo de rotas autenticadas
    authenticated := r.Group("")
    authenticated.Use(middleware.AuthRequired())
    {
        // Rotas de Produtos
        produtos := authenticated.Group("/produtos")
        {
            produtos.GET("", handlers.GetProdutos)
            produtos.GET("/:id", handlers.GetProduto)
            produtos.POST("", middleware.ManagerRequired(), handlers.CreateProduto)
            produtos.PUT("/:id", middleware.ManagerRequired(), handlers.UpdateProduto)
            produtos.DELETE("/:id", middleware.AdminRequired(), handlers.DeleteProduto)
            
            produtos.GET("/categoria/:categoria", handlers.GetProdutosPorCategoria)
            produtos.GET("/busca", handlers.BuscarProdutos)
            produtos.PATCH("/:id/estoque", middleware.ManagerRequired(), handlers.AtualizarEstoque)
            produtos.PATCH("/:id/preco", middleware.ManagerRequired(), handlers.AtualizarPreco)
            produtos.GET("/baixo-estoque", handlers.GetProdutosBaixoEstoque)
            produtos.POST("/:id/imagem", middleware.ManagerRequired(), handlers.UploadImagemProduto)
        }

        // Rotas de Relatórios (apenas admin e manager)
        relatorios := authenticated.Group("/relatorios")
        relatorios.Use(middleware.ManagerRequired())
        {
            relatorios.GET("/estoque", handlers.RelatorioEstoque)
            relatorios.GET("/produtos-mais-vendidos", handlers.RelatorioProdutosMaisVendidos)
            relatorios.GET("/valor-total-estoque", handlers.RelatorioValorTotalEstoque)
        }
    }

    r.Run(":8080")
}
