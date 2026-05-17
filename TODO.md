#  **TODO**

## **Backend**

- [x] Admin routes
  - [x] PATCH /users/{id} <!--Username-->
  - [x] PATCH /users/{id}/reset-password
  - [x] DELETE /users/{id}
- [x] User routes
  - [x] PATCH /change-username <!--Username-->
  - [x] PATCH /reset-password
- [ ] Initial admin seed script
- [ ] Admin deletion
- [x] Pagination
  - [x] TLDRs
  - [x] Users (Admin)
- [ ] Batch deletion of TLDRs
  - [ ] Users
  - [ ] Admin
- [x] AI-generated titles in TLDR create flow
- [ ] Dev endpoints in admin (DELETE /tldrs, DELETE /users)
- [x] API error handling in /summarize
- [ ] Flag for sensitive TLDRs
- [x] Id + CreatedAt to encode PageCursor to Base64 string
- [x] Refactor code to stay within 80-char line limit
- [ ] Documentation
- [ ] Testing
  - [x] config
  - [x] auth
  - [ ] core (stashed)

### **Future**

Plans to consider or implement in the future, if the app ever grows.

- Caching
- Indexing
- Role-Based Access Control
- Migrate to PostgreSQL
- TLDR Sharing (Public, Private, Shared)

## **Frontend**

Can start working after backend is polished.
