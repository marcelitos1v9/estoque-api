package handlers

import (
    "bytes"
    "context"
    "estoque-api/models"
    "estoque-api/auth"
    "estoque-api/database"
    "fmt"
    "io"
    "net/http"
    "time"
    "golang.org/x/crypto/bcrypt"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

// InitializeAuthHandlers inicializa as collections necessárias
func InitializeAuthHandlers() {
    userCollection = database.DB.Collection("users")
}

func Register(c *gin.Context) {
    // Criar uma estrutura específica para o registro
    var registerData struct {
        Nome     string `json:"nome" binding:"required"`
        Email    string `json:"email" binding:"required"`
        Senha    string `json:"senha" binding:"required"`
        Role     string `json:"role"`
    }

    // Log do body recebido
    body, _ := io.ReadAll(c.Request.Body)
    fmt.Printf("Body recebido: %s\n", string(body))
    // Restaura o body para poder ser lido novamente
    c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

    if err := c.ShouldBindJSON(&registerData); err != nil {
        fmt.Printf("Erro no binding: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao ler dados do usuário", "details": err.Error()})
        return
    }

    // Log para debug
    fmt.Printf("Dados após binding: %+v\n", registerData)

    // Validações básicas
    if registerData.Email == "" || registerData.Senha == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email e senha são obrigatórios"})
        return
    }

    // Verifica se o email já existe
    var existingUser models.User
    err := userCollection.FindOne(context.Background(), bson.M{"email": registerData.Email}).Decode(&existingUser)
    if err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email já cadastrado"})
        return
    }

    // Hash da senha
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Senha), bcrypt.DefaultCost)
    if err != nil {
        fmt.Printf("Erro ao gerar hash da senha: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
        return
    }

    // Criar novo usuário
    newUser := models.User{
        ID:           primitive.NewObjectID(),
        Nome:         registerData.Nome,
        Email:        registerData.Email,
        Senha:        string(hashedPassword),
        Role:         registerData.Role,
        DataCriacao:  time.Now(),
        Ativo:        true,
    }

    if newUser.Role == "" {
        newUser.Role = "user" // Role padrão
    }

    // Log para debug
    fmt.Printf("Dados do usuário antes de salvar: %+v\n", newUser)

    _, err = userCollection.InsertOne(context.Background(), newUser)
    if err != nil {
        fmt.Printf("Erro ao salvar usuário: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
        return
    }

    // Remove a senha antes de enviar a resposta
    newUser.Senha = ""
    c.JSON(http.StatusCreated, newUser)
}

func Login(c *gin.Context) {
    var credentials struct {
        Email    string `json:"email"`
        Senha    string `json:"senha"`
    }

    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao ler credenciais", "details": err.Error()})
        return
    }

    // Validações básicas
    if credentials.Email == "" || credentials.Senha == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email e senha são obrigatórios"})
        return
    }

    // Log para debug
    fmt.Printf("Tentativa de login para email: %s\n", credentials.Email)
    fmt.Printf("Senha fornecida: %s\n", credentials.Senha)

    var user models.User
    err := userCollection.FindOne(context.Background(), bson.M{"email": credentials.Email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            fmt.Printf("Usuário não encontrado para email: %s\n", credentials.Email)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Email não encontrado"})
            return
        }
        fmt.Printf("Erro ao buscar usuário: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário"})
        return
    }

    // Log para debug
    fmt.Printf("Usuário encontrado: %+v\n", user)
    fmt.Printf("Hash da senha armazenada: %s\n", user.Senha)

    // Comparação das senhas
    err = bcrypt.CompareHashAndPassword([]byte(user.Senha), []byte(credentials.Senha))
    if err != nil {
        fmt.Printf("Erro na comparação de senhas: %v\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Senha incorreta"})
        return
    }

    // Se chegou aqui, a senha está correta
    token, err := auth.GenerateToken(user.ID.Hex(), user.Role)
    if err != nil {
        fmt.Printf("Erro ao gerar token: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
        return
    }

    // Atualiza último acesso
    _, err = userCollection.UpdateOne(
        context.Background(),
        bson.M{"_id": user.ID},
        bson.M{"$set": bson.M{"ultimo_acesso": time.Now()}},
    )

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id": user.ID,
            "nome": user.Nome,
            "email": user.Email,
            "role": user.Role,
        },
    })
}