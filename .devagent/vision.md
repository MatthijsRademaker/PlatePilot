# Project Vision

> This document provides long-term context for the Product Owner agent to generate
> aligned and strategic user stories. Edit this file to guide the autonomous
> development of your project.

## Project Overview

**PlatePilot** is an intelligent meal planning and healthy lifestyle companion app with gamification elements. The app is built around the concept of multiple AI agent personas that guide users through different aspects of healthy living. Each agent specializes in a specific domain and operates as a distinct section/mini-app within the larger application.

The core problem PlatePilot solves: Making healthy eating accessible, personalized, and engaging by combining AI-powered recommendations with an intuitive, agent-guided experience.

## Target Users

- Health-conscious individuals looking for meal planning assistance
- People who want personalized recipe recommendations based on their preferences
- Users who benefit from guided, AI-assisted experiences rather than complex manual tools
- Anyone seeking to build sustainable healthy eating habits through gamification

## MVP Scope

The MVP focuses on the **Mealplanner Agent** - the first of three planned AI personas.

- [ ] Mealplanner Agent with vector-based recipe suggestions
- [ ] Recipe browsing and search (semantic similarity using pgvector)
- [ ] Weekly meal plan creation and management
- [ ] Basic recipe detail views with ingredients and instructions
- [ ] Mobile-friendly UI with intuitive navigation

## Short-Term Goals (Next Sprint)

- Implement vector-based recipe similarity search in the frontend
- Create the meal planning UI (week view with meal slots)
- Connect frontend to existing backend APIs (BFF → gRPC services)
- Add basic recipe filtering (by cuisine, ingredients, allergies)

## Medium-Term Goals (1-3 Months)

- **Calorie Tracker Agent**: Track daily intake, work with mealplanner to create balanced weekly plans
- **Recipe Creator Agent**: Generate new recipes based on user preferences and available ingredients
- Integrate real AI embeddings (Azure OpenAI) for semantic recipe search
- Add gamification elements (streaks, achievements, points for healthy choices)
- User profiles with dietary preferences and restrictions

## Long-Term Vision (6+ Months)

The ultimate goal is a **multi-agent collaborative system** where:
- All three agents communicate and coordinate to assist the user
- The Mealplanner consults the Calorie Tracker for nutritional balance
- The Recipe Creator suggests meals that fit the user's weekly plan
- Agents proactively offer suggestions based on user patterns
- Rich gamification: challenges, social features, progress tracking
- Voice/conversational interface with agent personas

## Technical Constraints

- **Backend**: Go microservices (Recipe API, MealPlanner API, Mobile BFF)
- **Frontend**: Vue.js 3 + Quasar 2 (vertical slice architecture)
- **Vector Search**: PostgreSQL with pgvector extension
- **Communication**: gRPC between services, REST for client-facing BFF
- **Messaging**: RabbitMQ for event-driven updates
- **No ORM**: Explicit SQL queries with pgx
- **Hobby project**: Simple solutions preferred, no enterprise patterns

## Out of Scope

- Social features (sharing, friends, leaderboards) - future consideration
- E-commerce/grocery ordering integration
- Wearable device integration
- Complex nutritional analysis beyond basic calorie tracking
- Multi-language support (English only for now)
- Authentication/authorization (not yet implemented)

## Quality Standards

- Code must have tests (table-driven tests for Go)
- Must follow existing code patterns (vertical slices in frontend, CQRS in backend)
- Performance considerations (vector search should be fast)
- Intuitive UX - the app should feel like talking to helpful assistants
- Mobile-first design with Quasar components

## Priority Guidelines

1. Security and stability
2. Core functionality (MVP: mealplanning with vector search)
3. User experience (intuitive navigation, agent personas)
4. Performance optimization (fast vector queries)
5. Gamification elements
6. Multi-agent coordination

## Design Assets

Place design files in the designs/ folder and reference them here:

- designs/agent-personas.md - Description of each AI agent persona
- designs/app-navigation.png - Navigation flow between agent sections
- designs/mealplan-week-view.png - Weekly meal planning interface

## Notes for Autonomous Development

- Focus on the Mealplanner Agent first - it's the MVP foundation
- Each agent should feel like a distinct mini-app with its own personality
- Keep the three agents loosely coupled initially, plan for integration later
- Prefer small, well-defined tasks over large features
- Consider dependencies between features
- The backend infrastructure (Go services, vector search) is already in place
- Frontend vertical slices align well with the agent/section concept
