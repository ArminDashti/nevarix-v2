import { NavLink, Route, Routes } from "react-router-dom";
import { routes } from "./routes";
import { cn } from "@/lib/utils";

export function AppShell() {
  return (
    <div className="flex min-h-screen">
      <aside className="w-56 border-r border-border bg-muted/30 p-4">
        <h1 className="mb-6 text-lg font-bold">Nevarix Manager</h1>
        <nav className="flex flex-col gap-1">
          {routes.map((route) => (
            <NavLink
              key={route.path}
              to={route.path}
              end={route.path === "/"}
              className={({ isActive }) =>
                cn(
                  "rounded-md px-3 py-2 text-sm hover:bg-muted",
                  isActive && "bg-primary text-primary-foreground hover:opacity-90",
                )
              }
            >
              {route.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <main className="flex-1 p-6">
        <Routes>
          {routes.map((route) => (
            <Route key={route.path} path={route.path} element={route.element} />
          ))}
        </Routes>
      </main>
    </div>
  );
}
