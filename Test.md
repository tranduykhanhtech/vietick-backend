# Tổng hợp lệnh curl test API

> **Lưu ý cho Windows CMD:**  
> Khi dùng curl với tham số `-d`, hãy dùng dấu nháy kép `"..."` và escape dấu nháy bên trong bằng `\"`.  
> Ví dụ:  
> ```sh
> curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"123456\",\"full_name\":\"Test User\"}"
> ```
> Nếu dùng Git Bash hoặc PowerShell, có thể dùng dấu nháy đơn `'...'` như bình thường.

---

## 1. Auth APIs

```sh
✅curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"123456\",\"full_name\":\"Test User\"}" /
curl -X POST https://vietick-backend.onrender.com/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"123456\"}"
curl -X POST http://localhost:8080/api/v1/auth/refresh -H "Content-Type: application/json" -d "{\"refresh_token\":\"{refresh_token}\"}"
curl -X POST http://localhost:8080/api/v1/auth/verify-email -H "Content-Type: application/json" -d "{\"token\":\"{verify_token}\"}"
curl -X POST http://localhost:8080/api/v1/auth/logout -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/auth/logout-all -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/auth/change-password -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"old_password\":\"123456\",\"new_password\":\"654321\"}"
curl -X GET http://localhost:8080/api/v1/auth/me -H "Authorization: Bearer {token}"
curl -X GET http://localhost:8080/api/v1/auth/check -H "Authorization: Bearer {token}"
```

## 2. User APIs

```sh
curl http://localhost:8080/api/v1/users/check-username?username=testuser
curl http://localhost:8080/api/v1/users/check-email?email=test@example.com
curl http://localhost:8080/api/v1/users/{id} -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/username/{username} -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/{id}/stats -H "Authorization: Bearer {token}"
curl -X PUT http://localhost:8080/api/v1/users/me -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"full_name\":\"New Name\"}"
curl -X PUT http://localhost:8080/api/v1/users/me/username -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"username\":\"newusername\"}"
curl -X PUT http://localhost:8080/api/v1/users/me/email -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"email\":\"new@example.com\"}"
curl http://localhost:8080/api/v1/users/search?q=test -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/recommended -H "Authorization: Bearer {token}"
```

## 3. Follow APIs

```sh
curl -X POST http://localhost:8080/api/v1/users/{id}/follow -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/users/{id}/unfollow -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/users/{id}/toggle-follow -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/{id}/follow-status -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/{id}/followers -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/{id}/following -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/{id}/follow-counts -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/{id}/mutual-follows -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/{id}/relationship -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/users/{id}/follow-stats -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/follows/bulk-follow -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"user_ids\":[1,2,3]}"
curl -X POST http://localhost:8080/api/v1/follows/bulk-unfollow -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"user_ids\":[1,2,3]}"
```

## 4. Post APIs

```sh
curl -X POST http://localhost:8080/api/v1/posts -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"content\":\"Hello world!\",\"image_urls\":[\"url1\",\"url2\"]}"
curl http://localhost:8080/api/v1/posts/feed -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/posts/explore -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/posts/search?q=hello -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/posts/{id} -H "Authorization: Bearer {token}"
curl -X PUT http://localhost:8080/api/v1/posts/{id} -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"content\":\"Updated content\"}"
curl -X DELETE http://localhost:8080/api/v1/posts/{id} -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/posts/{id}/stats -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/posts/{id}/like -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/posts/{id}/unlike -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/posts/{id}/toggle-like -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/posts/user/{user_id} -H "Authorization: Bearer {token}"
```

## 5. Comment APIs

```sh
curl -X POST http://localhost:8080/api/v1/posts/{id}/comments -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"content\":\"Nice post!\"}"
curl http://localhost:8080/api/v1/posts/{id}/comments -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/comments/{id} -H "Authorization: Bearer {token}"
curl -X PUT http://localhost:8080/api/v1/comments/{id} -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"content\":\"Updated comment\"}"
curl -X DELETE http://localhost:8080/api/v1/comments/{id} -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/comments/{id}/stats -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/comments/{id}/like -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/comments/{id}/unlike -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/comments/{id}/toggle-like -H "Authorization: Bearer {token}"
```

## 6. Verification APIs

```sh
curl http://localhost:8080/api/v1/verification/requirements
curl http://localhost:8080/api/v1/verification/me -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/verification/can-submit -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/verification/verified-users -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/verification/submit -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"full_name\":\"Test\",\"id_number\":\"123456\",\"id_type\":\"CCCD\",\"front_image_url\":\"url1\",\"selfie_image_url\":\"url2\"}"
curl http://localhost:8080/api/v1/verification/pending -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/verification/all -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/verification/stats -H "Authorization: Bearer {token}"
curl http://localhost:8080/api/v1/verification/{id} -H "Authorization: Bearer {token}"
curl -X POST http://localhost:8080/api/v1/verification/{id}/review -H "Authorization: Bearer {token}" -H "Content-Type: application/json" -d "{\"status\":\"approved\",\"admin_notes\":\"OK\"}"
curl -X DELETE http://localhost:8080/api/v1/verification/{id} -H "Authorization: Bearer {token}"
```

## 7. Health check

```sh
curl http://localhost:8080/health
```

> **Lưu ý:** Thay `{token}`, `{id}`, `{user_id}`, `{username}`, `{post_id}`, `{comment_id}`... bằng giá trị thực tế.