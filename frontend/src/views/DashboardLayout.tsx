import React from 'react';
import { Outlet, Link, useNavigate, useLocation } from '@tanstack/react-router';
import { useAuth } from '../context/AuthContext';
import { LayoutDashboard, BookOpen, LogOut, ShieldAlert, Cpu, Server as ServerIcon, Clock } from 'lucide-react';

export const DashboardLayout: React.FC = () => {
  const { user, isAuthenticated, isLoading, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  React.useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate({ to: '/login' });
    }
  }, [isLoading, isAuthenticated, navigate]);

  const handleLogout = () => {
    logout();
    navigate({ to: '/login' });
  };

  if (isLoading) {
    return (
      <div style={{ display: 'flex', height: '100vh', alignItems: 'center', justifyContent: 'center', backgroundColor: 'var(--bg-primary)' }}>
        <div style={{ textAlign: 'center' }}>
          <div style={{ width: '40px', height: '40px', border: '3px solid var(--border)', borderTopColor: 'var(--accent)', borderRadius: '50%', animation: 'spin 1s linear infinite', margin: '0 auto 16px' }} />
          <p style={{ color: 'var(--text-secondary)' }}>Loading dashboard...</p>
        </div>
        <style>{`
          @keyframes spin {
            to { transform: rotate(360deg); }
          }
        `}</style>
      </div>
    );
  }

  if (!user) return null;

  // Derive header title based on current route path
  let headerTitle = 'IT Systems Overview';
  if (location.pathname.includes('/systems/')) {
    headerTitle = 'System Management Details';
  } else if (location.pathname.includes('/servers')) {
    headerTitle = 'Host Servers Pool';
  } else if (location.pathname.includes('/documentations/new')) {
    headerTitle = 'Create Global System Documentation';
  } else if (location.pathname.includes('/pending-requests')) {
    headerTitle = 'Pending Feature Requests';
  }

  return (
    <div className="dashboard-shell">
      {/* Sidebar navigation */}
      <aside className="sidebar">
        <div className="sidebar-header">
          <div className="sidebar-logo">
            <Cpu size={18} />
          </div>
          <span className="sidebar-brand">IT DASHBOARD</span>
        </div>

        <nav className="sidebar-nav">
          <Link
            to="/dashboard"
            className={`sidebar-item ${location.pathname === '/dashboard' || location.pathname === '/dashboard/' ? 'active' : ''}`}
          >
            <LayoutDashboard size={18} />
            <span>Systems Dashboard</span>
          </Link>

          <Link
            to="/dashboard/servers"
            className={`sidebar-item ${location.pathname.includes('/servers') ? 'active' : ''}`}
          >
            <ServerIcon size={18} />
            <span>Servers Manager</span>
          </Link>

          <Link
            to="/dashboard/documentations/new"
            className={`sidebar-item ${location.pathname.includes('/documentations/new') ? 'active' : ''}`}
          >
            <BookOpen size={18} />
            <span>Add Documentation</span>
          </Link>

          <Link
            to="/dashboard/pending-requests"
            className={`sidebar-item ${location.pathname.includes('/pending-requests') ? 'active' : ''}`}
          >
            <Clock size={18} />
            <span>Pending Requests</span>
          </Link>
        </nav>

        <div className="sidebar-footer">
          <div className="sidebar-user">
            <div className="user-avatar">
              {user.username.substring(0, 2).toUpperCase()}
            </div>
            <div className="user-info">
              <div className="user-name">{user.username}</div>
              <div className="user-role">{user.role}</div>
            </div>
          </div>

          <button
            onClick={handleLogout}
            className="btn btn-secondary btn-sm"
            style={{ width: '100%', display: 'flex', justifyContent: 'center', gap: '8px', color: 'var(--danger)' }}
          >
            <LogOut size={14} />
            <span>Logout</span>
          </button>
        </div>
      </aside>

      {/* Main dashboard content */}
      <div className="main-content">
        <header className="header">
          <h1 className="header-title">{headerTitle}</h1>
          <div className="header-actions">
            <span style={{ fontSize: '13px', display: 'flex', alignItems: 'center', gap: '6px', color: 'var(--text-secondary)' }}>
              <ShieldAlert size={14} style={{ color: user.role === 'Viewer' ? 'var(--warning)' : 'var(--success)' }} />
              Role: <strong>{user.role}</strong>
            </span>
          </div>
        </header>

        <main className="content-body">
          <Outlet />
        </main>
      </div>
    </div>
  );
};
