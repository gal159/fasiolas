# 🔐 ADMIN ACCOUNT CREATED

**Date**: January 25, 2026  
**Status**: ✅ **ACTIVE**

---

## 👤 Admin Account Details

### Primary Admin
```
Username: Džiugas Elite
Email:    ciuta.dziugas@gmail.com
Role:     admin
User ID:  1
Status:   ✅ ACTIVE
```

**Password**: Use your existing Google OAuth login or your account password

---

## 🚀 How to Login as Admin

### Option 1: Via Frontend (Browser)
1. Open: `http://localhost:3000`
2. Click "Login with Google"
3. Use email: `ciuta.dziugas@gmail.com`
4. You now have **admin privileges**!

### Option 2: Via API
```bash
# Get your JWT token with admin role
curl -X GET "http://localhost:8080/api/v1/auth/google" \
  -H "Content-Type: application/json"

# Follow the OAuth flow to get your admin token
```

---

## 🎯 Admin Capabilities

Now you can:

### ✅ User Management
```bash
# List all users (replace TOKEN with your JWT)
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer YOUR_TOKEN"

# View specific user
curl http://localhost:8080/api/v1/admin/users/2 \
  -H "Authorization: Bearer YOUR_TOKEN"

# Change user role to spectator
curl -X PUT http://localhost:8080/api/v1/admin/users/2/role \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "spectator"}'

# Change user role back to player
curl -X PUT http://localhost:8080/api/v1/admin/users/2/role \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "player"}'

# 🆕 DELETE USER (PERMANENT!)
curl -X DELETE http://localhost:8080/api/v1/admin/users/2 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### ✅ Admin Panel UI (NEW!)
**Access the admin panel in your browser:**
1. Login to `http://localhost:3000`
2. Click **"Admin Panel"** in the navigation (only visible to admins)
3. View all users in a table format
4. Change user roles with a dropdown
5. Delete users with the **🗑️ Delete** button
6. See real-time statistics (total users, players, spectators, admins)

### ✅ System Statistics
```bash
# Get platform statistics
curl http://localhost:8080/api/v1/admin/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### ✅ All Player Capabilities
- Create games
- Join games
- Play cards
- Draw cards
- Skip turns
- Call cheat

---

## 📊 Current User Roster

| ID | Username | Email | Role |
|----|----------|-------|------|
| 1 | Džiugas Elite | ciuta.dziugas@gmail.com | **admin** |
| 2 | Rasa Ciutiene | ciutiene.rasa@gmail.com | player |

---

## 🔧 Additional Admin Commands

### Create Another Admin (if needed)
```bash
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game \
  -c "UPDATE users SET role = 'admin' WHERE email = 'another@email.com';"
```

### Change User to Spectator
```bash
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game \
  -c "UPDATE users SET role = 'spectator' WHERE id = 2;"
```

### View All Users
```bash
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game \
  -c "SELECT id, username, email, role FROM users ORDER BY id;"
```

---

## 🎮 Testing Your Admin Access

### Step 1: Login
1. Go to `http://localhost:3000`
2. Login with Google using `ciuta.dziugas@gmail.com`
3. Your admin status is automatically detected

### Step 2: Test Admin API
```bash
# First, get your token by logging in through the frontend
# Then test an admin endpoint:

export ADMIN_TOKEN="your_jwt_token_here"

# List users
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### Step 3: Test Role Management
```bash
# Change Rasa's role to spectator
curl -X PUT http://localhost:8080/api/v1/admin/users/2/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "spectator"}'

# Verify the change
curl http://localhost:8080/api/v1/admin/users/2 \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Change back to player
curl -X PUT http://localhost:8080/api/v1/admin/users/2/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "player"}'
```

---

## ⚡ Quick Commands Reference

```bash
# View all users with roles
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game \
  -c "SELECT id, username, email, role FROM users;"

# Make user admin
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game \
  -c "UPDATE users SET role = 'admin' WHERE id = USER_ID;"

# Make user spectator
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game \
  -c "UPDATE users SET role = 'spectator' WHERE id = USER_ID;"

# Make user player
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game \
  -c "UPDATE users SET role = 'player' WHERE id = USER_ID;"
```

---

## 🎯 What's Next?

1. **Login** to the frontend as admin
2. **Test** admin API endpoints
3. **Manage** user roles as needed
4. **Create** games and play (you still have player abilities!)
5. **Monitor** platform with admin stats

---

## 📚 More Information

- Full documentation: `FINAL_PROJECT_REPORT.md`
- Technical details: `ROLE_BASED_ACCESS_CONTROL.md`
- Quick reference: `ROLE_SYSTEM_QUICK_REF.md`
- Troubleshooting: `ROLE_SYSTEM_TROUBLESHOOTING.md`

---

## ✅ Summary

**Admin Account Created**: ✅  
**Username**: Džiugas Elite  
**Email**: ciuta.dziugas@gmail.com  
**Role**: admin  
**Status**: Active and ready to use!

**Next Step**: Login to `http://localhost:3000` and start managing the platform!

---

**Created**: January 25, 2026  
**Status**: ✅ ACTIVE  
