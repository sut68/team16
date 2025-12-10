package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/config"
	"backend/entity"
	"backend/services"
)
// helper to get user id from gin context
func getUserIDFromContext(c *gin.Context) (uint, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return uint(t), true
	case int:
		return uint(t), true
	case int64:
		return uint(t), true
	case uint:
		return t, true
	case uint64:
		return uint(t), true
	case string:
		n, err := strconv.ParseUint(t, 10, 64)
		if err == nil {
			return uint(n), true
		}
	}
	return 0, false
}

type createUserPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type updateUserPayload struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Role     *string `json:"role,omitempty"`
}

func ListUsers(c *gin.Context) {
	var users []entity.User
	if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

type CreateUserPayload struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
  Role     string `json:"role" binding:"required"` // expect "admin" or "student" or "user"
}

func CreateUser(c *gin.Context) {
  // get creator id from context (AuthMiddleware should set)
  creator, exists := c.Get("user_id")
  if !exists {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
    return
  }
  creatorID := creator.(uint)
  if creatorID != 1 { // only admin id=1 may create admins or manage users per your rule
    c.JSON(http.StatusForbidden, gin.H{"error": "only admin id=1 can create users"})
    return
  }

  var payload CreateUserPayload
  if err := c.ShouldBindJSON(&payload); err != nil {
    log.Printf("CreateUser bind error: %v", err)
    // return detailed error message temporarily for debugging
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "detail": err.Error()})
    return
  }

  // Normalize role
  roleName := strings.ToLower(strings.TrimSpace(payload.Role))
  var role entity.Role
  if err := config.DB.Where("LOWER(name) = ?", roleName).First(&role).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
    return
  }

  hashed, err := services.HashPassword(payload.Password)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
    return
  }

  user := entity.User{
    Username: payload.Username,
    Password: hashed,
    RoleID:   &role.ID,
  }

  if err := config.DB.Create(&user).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed create user", "detail": err.Error()})
    return
  }

  c.JSON(http.StatusCreated, gin.H{"message": "created", "id": user.ID})
}


func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	targetID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	targetID := uint(targetID64)

	callerID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var payload updateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var target entity.User
	if err := config.DB.Preload("Role").First(&target, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup user"})
		return
	}

	// load caller
	var caller entity.User
	if err := config.DB.Preload("Role").First(&caller, callerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load caller"})
		return
	}
	callerRole := ""
	if caller.Role != nil {
		callerRole = strings.ToLower(caller.Role.Name)
	}

	// if change role requested
	if payload.Role != nil {
		newRoleName := strings.ToLower(strings.TrimSpace(*payload.Role))
		if newRoleName == "admin" && callerID != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admin id=1 can set Admin role"})
			return
		}
		if newRoleName == "student" && !(callerID == 1 || callerRole == "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admin can set Student role"})
			return
		}
		var newRole entity.Role
		if err := config.DB.Where("LOWER(name)=?", newRoleName).First(&newRole).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
		target.RoleID = &newRole.ID
	}

	if payload.Username != nil {
		target.Username = *payload.Username
	}
	if payload.Password != nil {
		hash, err := services.HashPassword(*payload.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		target.Password = hash
	}

	if err := config.DB.Save(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	targetID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	targetID := uint(targetID64)

	callerID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var target entity.User
	if err := config.DB.Preload("Role").First(&target, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to lookup user"})
		return
	}

	// cannot delete self
	if callerID == targetID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete yourself"})
		return
	}

	targetRoleName := ""
	if target.Role != nil {
		targetRoleName = strings.ToLower(target.Role.Name)
	}

	if targetRoleName == "admin" && callerID != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin id=1 can delete admin accounts"})
		return
	}

	// only admin family can delete users
	var caller entity.User
	if err := config.DB.Preload("Role").First(&caller, callerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load caller"})
		return
	}
	callerRole := ""
	if caller.Role != nil {
		callerRole = strings.ToLower(caller.Role.Name)
	}
	if !(callerID == 1 || callerRole == "admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin users can delete accounts"})
		return
	}

	if err := config.DB.Delete(&entity.User{}, targetID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
