import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from '@tanstack/react-router';
import { api } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { DataTable } from '../components/DataTable';
import { Plus, Edit2, Trash2, Eye, X, ExternalLink } from 'lucide-react';
import type { ColumnDef } from '@tanstack/react-table';

interface Server {
  id: number;
  name: string;
  ip: string;
}

interface System {
  id: number;
  name: string;
  type: string;
  links: string;
  server_id: number;
  status: string;
  description: string;
  created_at: string;
}

export const SystemsOverview: React.FC = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingSystem, setEditingSystem] = useState<System | null>(null);
  
  // Form state fields aligned to backend CreateSystemRequest
  const [name, setName] = useState('');
  const [type, setType] = useState('');
  const [links, setLinks] = useState('');
  const [serverId, setServerId] = useState<string>('');
  const [status, setStatus] = useState('Active');
  const [description, setDescription] = useState('');
  const [error, setError] = useState<string | null>(null);

  // Fetch systems
  const { data: systems = [], isLoading: isLoadingSystems } = useQuery<System[]>({
    queryKey: ['systems'],
    queryFn: async () => {
      const response = await api.get('/api/systems/');
      return response.data;
    },
  });

  // Fetch servers to populate selection dropdown
  const { data: servers = [], isLoading: isLoadingServers } = useQuery<Server[]>({
    queryKey: ['servers'],
    queryFn: async () => {
      const response = await api.get('/api/servers/');
      return response.data;
    },
  });

  // Mutations
  const createMutation = useMutation({
    mutationFn: async (newSystem: any) => {
      const response = await api.post('/api/systems/', newSystem);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['systems'] });
      closeModal();
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || err.response?.data?.details || 'Failed to create system');
    }
  });

  const updateMutation = useMutation({
    mutationFn: async (updatedSystem: System) => {
      const response = await api.put(`/api/systems/${updatedSystem.id}`, {
        name: updatedSystem.name,
        type: updatedSystem.type,
        links: updatedSystem.links,
        server_id: Number(updatedSystem.server_id),
        status: updatedSystem.status,
        description: updatedSystem.description,
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['systems'] });
      closeModal();
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || err.response?.data?.details || 'Failed to update system');
    }
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/api/systems/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['systems'] });
    },
    onError: (err: any) => {
      alert(err.response?.data?.error || 'Failed to delete system');
    }
  });

  const openCreateModal = () => {
    setEditingSystem(null);
    setName('');
    setType('');
    setLinks('');
    setServerId(servers.length > 0 ? String(servers[0].id) : '');
    setStatus('Active');
    setDescription('');
    setError(null);
    setIsModalOpen(true);
  };

  const openEditModal = (system: System) => {
    setEditingSystem(system);
    setName(system.name);
    setType(system.type);
    setLinks(system.links);
    setServerId(String(system.server_id));
    setStatus(system.status);
    setDescription(system.description);
    setError(null);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setEditingSystem(null);
    setName('');
    setType('');
    setLinks('');
    setServerId('');
    setStatus('Active');
    setDescription('');
    setError(null);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !type.trim() || !links.trim() || !serverId || !status.trim() || !description.trim()) {
      setError('All fields are required');
      return;
    }

    const payload = {
      name,
      type,
      links,
      server_id: Number(serverId),
      status,
      description,
    };

    if (editingSystem) {
      updateMutation.mutate({
        ...editingSystem,
        ...payload,
      });
    } else {
      createMutation.mutate(payload);
    }
  };

  const handleDelete = (id: number, name: string) => {
    if (confirm(`Are you sure you want to delete the system "${name}"? This will delete all linked documentation, notes, and feature requests.`)) {
      deleteMutation.mutate(id);
    }
  };

  const isViewer = user?.role === 'Viewer';

  // Table Columns Definition
  const columns = React.useMemo<ColumnDef<System>[]>(() => [
    {
      accessorKey: 'name',
      header: 'System Name',
      cell: info => <span style={{ fontWeight: 600 }}>{info.getValue() as string}</span>,
    },
    {
      accessorKey: 'type',
      header: 'Type / Tech Stack',
      cell: info => <code style={{ background: 'var(--bg-tertiary)', padding: '2px 6px', borderRadius: '4px' }}>{info.getValue() as string}</code>,
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: info => {
        const val = info.getValue() as string;
        let badgeClass = 'badge-info';
        if (val === 'Active' || val === 'Online') badgeClass = 'badge-success';
        if (val === 'Maintenance') badgeClass = 'badge-warning';
        if (val === 'Offline' || val === 'Deprecated') badgeClass = 'badge-danger';
        return <span className={`badge ${badgeClass}`}>{val}</span>;
      },
    },
    {
      accessorKey: 'links',
      header: 'Links',
      cell: info => {
        const link = info.getValue() as string;
        return (
          <a href={link} target="_blank" rel="noreferrer" style={{ display: 'inline-flex', alignItems: 'center', gap: '4px' }}>
            <span>Link</span>
            <ExternalLink size={12} />
          </a>
        );
      },
    },
    {
      accessorKey: 'server_id',
      header: 'Server ID',
      cell: info => {
        const sId = info.getValue() as number;
        const matchingServer = servers.find(s => s.id === sId);
        return <span>{matchingServer ? `${matchingServer.name} (${matchingServer.ip})` : `Server ID: ${sId}`}</span>;
      },
    },
    {
      id: 'actions',
      header: 'Actions',
      cell: ({ row }) => {
        const system = row.original;
        return (
          <div style={{ display: 'flex', gap: '8px' }}>
            <Link
              to="/dashboard/systems/$systemId"
              params={{ systemId: String(system.id) }}
              className="btn btn-secondary btn-sm"
              style={{ padding: '6px' }}
            >
              <Eye size={14} />
            </Link>
            {!isViewer && (
              <>
                <button
                  onClick={() => openEditModal(system)}
                  className="btn btn-secondary btn-sm"
                  style={{ padding: '6px' }}
                >
                  <Edit2 size={14} />
                </button>
                <button
                  onClick={() => handleDelete(system.id, system.name)}
                  className="btn btn-danger btn-sm"
                  style={{ padding: '6px' }}
                >
                  <Trash2 size={14} />
                </button>
              </>
            )}
          </div>
        );
      },
    },
  ], [servers, isViewer]);

  if (isLoadingSystems || isLoadingServers) {
    return (
      <div style={{ textAlign: 'center', padding: '64px' }}>
        <div style={{ width: '45px', height: '45px', border: '3px solid var(--border)', borderTopColor: 'var(--accent)', borderRadius: '50%', animation: 'spin 1s linear infinite', margin: '0 auto 16px' }} />
        <p style={{ color: 'var(--text-secondary)' }}>Loading IT Dashboard metadata...</p>
      </div>
    );
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div>
          <p style={{ color: 'var(--text-secondary)' }}>Manage your enterprise application index, documentation, servers, and change requests.</p>
        </div>
        {!isViewer && (
          <button onClick={openCreateModal} className="btn btn-primary" disabled={servers.length === 0}>
            <Plus size={16} />
            <span>Add System</span>
          </button>
        )}
      </div>

      {servers.length === 0 && !isViewer && (
        <div className="login-error" style={{ marginBottom: '24px', background: 'var(--warning-glow)', borderColor: 'var(--warning)' }}>
          <span style={{ color: 'var(--warning)' }}>Warning: No servers are defined. Please go to the <strong>Servers Manager</strong> to add a server before creating a system.</span>
        </div>
      )}

      <div className="card" style={{ padding: '20px' }}>
        <DataTable columns={columns} data={systems} searchPlaceholder="Search systems by name, type, status..." />
      </div>

      {/* Create / Edit Modal */}
      {isModalOpen && (
        <div className="modal-overlay">
          <div className="modal-content">
            <div className="modal-header">
              <h3 className="modal-title">{editingSystem ? 'Edit System' : 'Create New System'}</h3>
              <button className="modal-close" onClick={closeModal}>
                <X size={18} />
              </button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="modal-body">
                {error && (
                  <div className="login-error" style={{ marginBottom: '16px' }}>
                    <span>{error}</span>
                  </div>
                )}
                <div className="form-group">
                  <label className="form-label">System Name</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. Payment Gateway Service"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Type / Tech Stack</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. Go + Gin, React + Vite, Postgres"
                    value={type}
                    onChange={(e) => setType(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Access Links / URL</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. https://api.prod.company.com"
                    value={links}
                    onChange={(e) => setLinks(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Linked Host Server</label>
                  <select
                    className="form-select"
                    value={serverId}
                    onChange={(e) => setServerId(e.target.value)}
                  >
                    <option value="" disabled>-- Select a server --</option>
                    {servers.map(s => (
                      <option key={s.id} value={s.id}>{s.name} ({s.ip})</option>
                    ))}
                  </select>
                </div>
                <div className="form-group">
                  <label className="form-label">Operation Status</label>
                  <select
                    className="form-select"
                    value={status}
                    onChange={(e) => setStatus(e.target.value)}
                  >
                    <option value="Active">Active</option>
                    <option value="Maintenance">Maintenance</option>
                    <option value="Deprecated">Deprecated</option>
                    <option value="Offline">Offline</option>
                  </select>
                </div>
                <div className="form-group">
                  <label className="form-label">System Description</label>
                  <textarea
                    className="form-textarea"
                    placeholder="Describe the system's role, owner team, critical business value, and dependencies."
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={closeModal}>Cancel</button>
                <button type="submit" className="btn btn-primary">
                  {editingSystem ? 'Save Changes' : 'Create System'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
