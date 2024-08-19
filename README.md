# PowerDefi Hackathon

## Introduction

This repository contains two codebases: the frontend and the backend. The frontend is written in Next.js and Typescript, while the backend is written in Go.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (version 1.16 or later)
- Node.js (version 20.15.1 or later)
- Make (optional, for using the Makefile)
- Postgresql database set up

## Installation

### Frontend Setup

Clone the repository: `git clone git@github.com:Powerdfi-com/powerdfi-hackathon.git`

1. Navigate to the frontend directory: `cd Frontend`
2. Install dependencies: `npm install`
3. Rename the `.env.example` file to `.env` and update the necessary environment variables.
4. Start the development server: `npm run dev`

### Backend Setup

1. Navigate to the backend directory: `cd Backend`
2. Install dependencies: `go mod tidy`
3. Rename the `.env.example` file to `.env` and update the necessary environment variables.
4. Run the database migrations using the Makefile command: `make migrate`
5. Start the backend server: `make run/api`

Please make sure to update the environment variables in the `.env` files for both the frontend and backend codebases before running them.
