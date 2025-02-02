package handlers

import (
    "context"
    "estoque-api/database"
    "estoque-api/models"
    "net/http"
    "path/filepath"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

// InitializeHandlers deve ser chamada após a conexão com o banco
func InitializeHandlers() {
    collection = database.DB.Collection("produtos")
}

func GetProdutos(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var produtos []models.Produto
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var produto models.Produto
        cursor.Decode(&produto)
        produtos = append(produtos, produto)
    }

    c.JSON(http.StatusOK, produtos)
}

func GetProduto(c *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(c.Param("id"))
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var produto models.Produto
    err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&produto)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
        return
    }

    c.JSON(http.StatusOK, produto)
}

func CreateProduto(c *gin.Context) {
    var produto models.Produto
    if err := c.ShouldBindJSON(&produto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    produto.ID = primitive.NewObjectID()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := collection.InsertOne(ctx, produto)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, produto)
}

func UpdateProduto(c *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(c.Param("id"))
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var produto models.Produto
    if err := c.ShouldBindJSON(&produto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    update := bson.M{"$set": bson.M{"nome": produto.Nome, "preco": produto.Preco, "estoque": produto.Estoque}}
    _, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Produto atualizado com sucesso"})
}

func DeleteProduto(c *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(c.Param("id"))
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Produto removido com sucesso"})
}

func GetProdutosPorCategoria(c *gin.Context) {
    categoria := c.Param("categoria")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var produtos []models.Produto
    cursor, err := collection.Find(ctx, bson.M{"categoria": categoria})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var produto models.Produto
        cursor.Decode(&produto)
        produtos = append(produtos, produto)
    }

    c.JSON(http.StatusOK, produtos)
}

func BuscarProdutos(c *gin.Context) {
    query := c.Query("q")
    filter := bson.M{
        "$or": []bson.M{
            {"nome": bson.M{"$regex": query, "$options": "i"}},
            {"descricao": bson.M{"$regex": query, "$options": "i"}},
            {"tags": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var produtos []models.Produto
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var produto models.Produto
        cursor.Decode(&produto)
        produtos = append(produtos, produto)
    }

    c.JSON(http.StatusOK, produtos)
}

func AtualizarEstoque(c *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(c.Param("id"))
    var dados struct {
        Quantidade int    `json:"quantidade"`
        Operacao   string `json:"operacao"` // "adicionar" ou "remover"
    }

    if err := c.ShouldBindJSON(&dados); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var updateQuery bson.M
    if dados.Operacao == "adicionar" {
        updateQuery = bson.M{"$inc": bson.M{"estoque": dados.Quantidade}}
    } else {
        updateQuery = bson.M{"$inc": bson.M{"estoque": -dados.Quantidade}}
    }

    result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, updateQuery)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Estoque atualizado", "modificados": result.ModifiedCount})
}

func GetProdutosBaixoEstoque(c *gin.Context) {
    limite, _ := strconv.Atoi(c.DefaultQuery("limite", "5"))
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var produtos []models.Produto
    cursor, err := collection.Find(ctx, bson.M{"estoque": bson.M{"$lt": limite}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var produto models.Produto
        cursor.Decode(&produto)
        produtos = append(produtos, produto)
    }

    c.JSON(http.StatusOK, produtos)
}

func AtualizarPreco(c *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(c.Param("id"))
    var dados struct {
        NovoPreco         float64 `json:"novo_preco"`
        PrecoPromocional  float64 `json:"preco_promocional,omitempty"`
    }

    if err := c.ShouldBindJSON(&dados); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    update := bson.M{
        "$set": bson.M{
            "preco": dados.NovoPreco,
            "preco_promocional": dados.PrecoPromocional,
            "ultima_atualizacao": time.Now(),
        },
    }

    result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Preço atualizado", "modificados": result.ModifiedCount})
}

func UploadImagemProduto(c *gin.Context) {
    id, _ := primitive.ObjectIDFromHex(c.Param("id"))
    
    // Recebe o arquivo
    file, err := c.FormFile("imagem")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não encontrado"})
        return
    }

    // Gera um nome único para o arquivo
    filename := id.Hex() + "_" + time.Now().Format("20060102150405") + filepath.Ext(file.Filename)
    
    // Salva o arquivo
    if err := c.SaveUploadedFile(file, "uploads/"+filename); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
        return
    }

    // Atualiza o produto com a URL da imagem
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    update := bson.M{
        "$set": bson.M{
            "imagem_url": "/uploads/" + filename,
            "ultima_atualizacao": time.Now(),
        },
    }

    _, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Imagem atualizada", "url": "/uploads/" + filename})
}

func RelatorioEstoque(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    pipeline := []bson.M{
        {
            "$group": bson.M{
                "_id": "$categoria",
                "total_produtos": bson.M{"$sum": 1},
                "total_estoque": bson.M{"$sum": "$estoque"},
                "valor_total": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$preco", "$estoque"}}},
            },
        },
    }

    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    var resultados []bson.M
    if err = cursor.All(ctx, &resultados); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resultados)
}

func RelatorioProdutosMaisVendidos(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    limite, _ := strconv.Atoi(c.DefaultQuery("limite", "10"))

    // Aqui você precisaria ter uma coleção de vendas para fazer essa análise
    // Este é um exemplo simplificado
    pipeline := []bson.M{
        {
            "$sort": bson.M{"vendas_totais": -1},
        },
        {
            "$limit": limite,
        },
    }

    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    var resultados []bson.M
    if err = cursor.All(ctx, &resultados); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resultados)
}

func RelatorioValorTotalEstoque(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    pipeline := []bson.M{
        {
            "$group": bson.M{
                "_id": nil,
                "valor_total": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$preco", "$estoque"}}},
                "total_itens": bson.M{"$sum": "$estoque"},
                "total_produtos": bson.M{"$sum": 1},
            },
        },
    }

    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    var resultado []bson.M
    if err = cursor.All(ctx, &resultado); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if len(resultado) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "valor_total": 0,
            "total_itens": 0,
            "total_produtos": 0,
        })
        return
    }

    c.JSON(http.StatusOK, resultado[0])
}
