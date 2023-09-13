# Bookings and Reservations

## Overview

A webpage that allows users to explore and book hotel rooms. This project was made in the context of an Udemy course to develop Go and Back-End skills.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Database](#database)
- [Testing](#testing)

## Features

What can be done in it?

- Check room availability
- Book rooms
- Cancel reservations

## Technologies

- Backend: Golang
- Frontend: HTML, CSS, JavaScript
- Database: PostgreSQL
- Other
    - [alex edwards SCS session management](https://github.com/alexedwards/scs)
    - [chi router](https://github.com/go-chi/chi)

## Database

The following diagram represents the relationships between the database tables:


```mermaid
erDiagram
    User ||--o{ Reservation : "Makes"
    Room ||--o{ Reservation : "Reserved for"
    Room ||--o{ RoomRestrictions : "Has"
    Restriction ||--o{ RoomRestrictions : "Applies to"
 ```

## Testing

Explain how to run tests for your project, including unit tests, integration tests, and end-to-end tests if applicable.
