# Unibox : The Unified Complaint Box for Universities

**Unibox** is a smart campus issue management system designed to bridge the communication gap between students and university administrations. It centralizes fragmented complaints—from leaky hostel pipes to mess food quality—and uses automated routing to ensure the right department sees the right issue at the right time.

---

## 🚀 Problem Statement
Universities often lack a centralized channel for infrastructure and service complaints. Students rely on physical registers, emails, or word-of-mouth, leading to:
- **Delayed Responses :** Complaints get lost in paperwork.
- **Accountability Gaps :** No way to track which department is lagging.
- **Redundancy :** Multiple students reporting the same issue.

## ✨ Features

### For Students (Mobile/Web)
- **Centralized Filing :** A single form for all campus issues (Hostel, Mess, Academics, Infrastructure).
- **Real-time Tracking :** A personal dashboard to monitor ticket status (Pending, In-Progress, Resolved, Rejected).
- **Evidence Upload :** Attach images to provide context to the administration.

### For Admins (Desktop-Optimized)
- **Smart Routing :** Automated delivery of complaints to specific department dashboards.
- **Status Management :** Update progress to keep students informed.
- **Manual Reroute (Fail-safe) :** If the system misroutes a ticket, admins can instantly transfer it to the correct department.
- **Role-Based Access :** Admins only see tickets relevant to their specific department.

---

## 🛠️ Tech Stack
- **Frontend :** Nijor + Tailwind
- **Backend :** Go + Fiber
- **Database :** PostgreSQL + Redis (Dual-table architecture for Student and Admin profiles)
- **Authentication :** JWT + Refresh Key Rotation

---

## 🏗️ System Architecture


1. **Submission :** Student submits a complaint via the `Student Profile`.
2. **Routing :** The backend logic analyzes the category and pushes the data to the corresponding `Admin Department`.
3. **Action :** Admin logs in via **Desktop** -> Views Queue -> Updates Status.
4. **Loop :** Student receives real-time updates on their dashboard.

---

## 📸 Screenshots (Coming Soon)
* [Student Dashboard View]
* [Admin Department Queue]
* [Complaint Filing Form]

---

## 🚦 Current Status & Roadmap

- [x] Student Authentication
- [x] Student Dashboard UI
- [ ] Smart Routing Logic (In Development)
- [ ] Admin Role-Based Dashboards
- [ ] Image Upload Integration
- [ ] Push Notifications for Status Changes

---