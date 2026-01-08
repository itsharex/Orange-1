# Orange - Project Management Tool

Orange is a modern, lightweight project management application built with **Go (Wails)** and **Vue 3**. It is designed to help small teams and freelancers manage projects, payments, and client relationships efficiently.

## Features

- **Project Management**: Create, track, and manage projects with ease.
- **Payment Tracking**: Record payments, track milestones, and visualize income trends.
- **Dashboard**: Get a bird's-eye view of your business with real-time statistics and charts.
- **Notifications**: Stay updated with a built-in notification system.
- **Cross-Platform**: Runs natively on macOS, Windows, and Linux.

## Tech Stack

- **Frontend**: Vue 3, TypeScript, Vite, Pinia
- **Backend**: Go, Wails v3, GORM (SQLite)
- **Styling**: Custom CSS (Glassmorphism design)

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- NPM or Yarn

### Development

1.  Clone the repository:

    ```bash
    git clone https://github.com/FruitsAI/Orange.git
    cd Orange
    ```

2.  Install frontend dependencies:

    ```bash
    cd frontend
    npm install
    cd ..
    ```

3.  Run the application in development mode:
    ```bash
    wails3 dev
    ```

### Build

To build the application for production:

```bash
wails3 build
```

The output binary will be located in the `build/` directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Copyright (c) 2026 FruitsAI
