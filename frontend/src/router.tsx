import { createRootRoute, createRoute, createRouter, Navigate } from '@tanstack/react-router';
import { LoginView } from './views/LoginView';
import { DashboardLayout } from './views/DashboardLayout';
import { SystemsOverview } from './views/SystemsOverview';
import { ServersOverview } from './views/ServersOverview';
import { SystemDetail } from './views/SystemDetail';
import { GlobalAddDocumentation } from './views/GlobalAddDocumentation';
import { PendingRequests } from './views/PendingRequests';
import { Outlet } from '@tanstack/react-router';

// 1. Root route
const rootRoute = createRootRoute({
  component: () => <Outlet />,
});

// 2. Index redirect route
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: () => <Navigate to="/dashboard" replace />,
});

// 3. Login route
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginView,
});

// 4. Dashboard layout route (Secured shell)
const dashboardRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/dashboard',
  component: DashboardLayout,
});

// 5. Dashboard index route (Lists IT systems)
const dashboardIndexRoute = createRoute({
  getParentRoute: () => dashboardRoute,
  path: '/',
  component: SystemsOverview,
});

// 6. Servers route (Lists host nodes)
const serversRoute = createRoute({
  getParentRoute: () => dashboardRoute,
  path: '/servers',
  component: ServersOverview,
});

// 7. System detail route (Manages specific systems with dynamic parameters)
const systemDetailRoute = createRoute({
  getParentRoute: () => dashboardRoute,
  path: '/systems/$systemId',
  component: SystemDetail,
});

// 8. Global add doc route (Directly posts writeups)
const createDocumentationRoute = createRoute({
  getParentRoute: () => dashboardRoute,
  path: '/documentations/new',
  component: GlobalAddDocumentation,
});

// 9. Pending feature requests route
const pendingRequestsRoute = createRoute({
  getParentRoute: () => dashboardRoute,
  path: '/pending-requests',
  component: PendingRequests,
});

// Build route tree
const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  dashboardRoute.addChildren([
    dashboardIndexRoute,
    serversRoute,
    systemDetailRoute,
    createDocumentationRoute,
    pendingRequestsRoute,
  ]),
]);

// Export router instance
export const router = createRouter({
  routeTree,
});

// Register router module type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
