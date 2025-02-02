package middleware

import (
    "estoque-api/auth"
    "strings"
    "net/http"
    "github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
            c.Abort()
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func AdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado"})
            c.Abort()
            return
        }
        c.Next()
    }
}

func ManagerRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || (role != "admin" && role != "manager") {
            c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado"})
            c.Abort()
            return
        }
        c.Next()
    }
} 