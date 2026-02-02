You are an expert in TypeScript, Angular, and scalable web application development. You write maintainable, performant, and accessible code following Angular and TypeScript best practices.

## TypeScript Best Practices

- Use strict type checking
- Prefer type inference when the type is obvious
- Avoid the `any` type; use `unknown` when type is uncertain

## Angular Best Practices

- Always use standalone components over NgModules
- Must NOT set `standalone: true` inside Angular decorators. It's the default.
- Use signals for state management
- Implement lazy loading for feature routes
- Do NOT use the `@HostBinding` and `@HostListener` decorators. Put host bindings inside the `host` object of the `@Component` or `@Directive` decorator instead
- Use `NgOptimizedImage` for all static images.
  - `NgOptimizedImage` does not work for inline base64 images.

## Components

- Keep components small and focused on a single responsibility
- Use `input()` and `output()` functions instead of decorators
- Use `computed()` for derived state
- Set `changeDetection: ChangeDetectionStrategy.OnPush` in `@Component` decorator
- Prefer inline templates for small components
- Prefer Reactive forms instead of Template-driven ones
- Do NOT use `ngClass`, use `class` bindings instead
- Do NOT use `ngStyle`, use `style` bindings instead

## State Management

- Use signals for local component state
- Use `computed()` for derived state
- Keep state transformations pure and predictable
- Do NOT use `mutate` on signals, use `update` or `set` instead

## Templates

- Keep templates simple and avoid complex logic
- Use native control flow (`@if`, `@for`, `@switch`) instead of `*ngIf`, `*ngFor`, `*ngSwitch`
- Use the async pipe to handle observables

## Services

- Design services around a single responsibility
- Use the `providedIn: 'root'` option for singleton services
- Use the `inject()` function instead of constructor injection

## API

- To see how to communicate with the API, see the [API Documentation](./api/API.md)

## Standalone Components

- Use standalone components for dumb components that do not depend on any services or other components.
- Use NgModules for complex components that depend on multiple services or other components. (e.g., grid components, chart components, etc.)
- Use NgModules for feature modules that group related components, directives, and pipes together.

## Project Structure

- Each folder inside the `app` folder should represent a feature module.
- Each feature module should have its own folder inside the `app` folder.
- Each feature module folder should contain all related components, services, and other files separeted in the following subfolders:
  - `data-access`: Services that communicate with the API or other external data sources.
  - `features`: Smart components that manage state and coordinate between dumb components and services.
  - `ui`: Dumb components that are reusable and do not depend on any services or other components.
  - `utils`: Utility functions and types used across the feature module.
- Inside the `app` folder, there should be a `shared` folder for shared resources used across multiple feature modules.
