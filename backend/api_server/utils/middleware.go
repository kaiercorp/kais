package utils

import (
	"errors"
	"net/http"
	"os"

	//todo: update version
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte(os.Getenv("SECRETKKAIER"))

// TokenData는 JWT에서 가져오는 사용자 정보를 저장하는 구조체입니다.
type TokenData struct {
	Username string `json:"username"` // 사용자 이름
	Group    int    `json:"group"`    // 사용자 그룹
}

// GetDataFromToken은 Gin 컨텍스트에서 JWT를 통해 추출된 사용자 정보를 반환합니다.
//
// 이 함수는 JWTAuthMiddleware에서 c.Set으로 저장된 "username" 및 "group" 값을 기반으로
// TokenData 구조체를 생성합니다.
//
// 사용자가 인증되지 않았거나, 타입이 올바르지 않은 경우 에러를 반환합니다.
//
// 반환값:
//   - *TokenData: 인증된 사용자 정보 (username, group)
//   - error: 인증 실패 또는 타입 오류 시 반환
func GetDataFromToken(c *gin.Context) (*TokenData, error) {
	username, exists := c.Get("username")
	if !exists {
		return nil, errors.New("user not authenticated")
	}

	group, exists := c.Get("group")
	if !exists {
		return nil, errors.New("group not found in token")
	}

	usernameStr, ok := username.(string)
	if !ok {
		return nil, errors.New("username is not a valid string")
	}

	groupInt, ok := group.(int)
	if !ok {
		return nil, errors.New("group is not a valid integer")
	}

	return &TokenData{
		Username: usernameStr,
		Group:    groupInt,
	}, nil
}

// JWTAuthMiddleware는 JWT 기반 인증을 수행하는 Gin 미들웨어입니다.
//
// Authorization 헤더로 전달된 Bearer 토큰을 파싱하여 유효성을 검증하고,
// 토큰에서 사용자 정보(username, group)를 추출하여 Gin 컨텍스트에 저장합니다.
//
// 이 미들웨어는 다음과 같은 경우 요청을 거부(401)합니다:
//   - Authorization 헤더가 없는 경우
//   - Bearer 형식이 아닌 경우
//   - 토큰이 만료되었거나 위조된 경우
//   - 토큰 payload에 필요한 정보가 없는 경우
//
// 컨텍스트 저장 항목:
//   - "username": string
//   - "group": int
//
// 반환값:
//   - gin.HandlerFunc: 인증 미들웨어
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := authHeader[len(bearerPrefix):]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		group, ok := claims["group"].(float64) // jwt는 숫자를 float64로 파싱함
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		// 이후 미들웨어나 핸들러에서 사용할 수 있도록 컨텍스트에 저장
		c.Set("username", username)
		c.Set("group", int(group))

		c.Next()
	}
}

// GroupMiddleware는 특정 그룹 사용자만 요청을 수행할 수 있도록 제한하는 미들웨어입니다.
//
// JWT 토큰에서 추출된 사용자 그룹 정보를 확인하여,
// allowedGroups에 포함되지 않은 경우 요청을 거부합니다 (403).
//
// 이 미들웨어는 다음과 같은 경우 요청을 차단합니다:
//   - JWT에서 사용자 정보 추출에 실패한 경우 (401)
//   - 허용되지 않은 그룹인 경우 (403)
//
// 파라미터:
//   - allowedGroups: 접근을 허용할 사용자 그룹의 정수 목록
//
// 반환값:
//   - gin.HandlerFunc: 권한 제어 미들웨어
func GroupMiddleware(allowedGroups ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, err := GetDataFromToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		isAllowed := false
		for _, allowedGroup := range allowedGroups {
			if userData.Group == allowedGroup {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
