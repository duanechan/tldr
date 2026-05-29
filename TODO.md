#  **TODO**

## **Backend**

- [x] Admin routes
  - [x] PATCH /users/{id} <!--Username-->
  - [x] PATCH /users/{id}/reset-password
  - [x] DELETE /users/{id}
- [x] User routes
  - [x] PATCH /change-username <!--Username-->
  - [x] PATCH /reset-password
- [x] Initial admin seed script
- [ ] Admin deletion
- [x] Pagination
  - [x] TLDRs
  - [x] Users (Admin)
- [x] Batch deletion of TLDRs
  - [x] Users
  - [x] Admin
- [x] AI-generated titles in TLDR create flow
- [x] Dev endpoints in admin (DELETE /tldrs, DELETE /users)
- [x] API error handling in /summarize
- [x] Flag for sensitive TLDRs
- [x] Id + CreatedAt to encode PageCursor to Base64 string
- [x] Refactor code to stay within 80-char line limit
- [ ] Documentation
  - [ ] README.md
  - [ ] Comments
- [x] Testing
  - [x] config
  - [x] auth
- [x] Refactor validate.String to return multiple errors
- [x] Add validate.String to register handler  
- [ ] Integration tests

### **Future**

Plans to consider or implement in the future, if the app ever grows.

- Caching
- Indexing
- Role-Based Access Control
- Migrate to PostgreSQL
- TLDR Sharing (Public, Private, Shared)
- Migrate to OAuth2

## **Frontend**

- [ ] Login
- [ ] Register
- [ ] TL;DRs in sidebar
- [ ] Admin side
