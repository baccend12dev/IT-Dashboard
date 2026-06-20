import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { DataTable } from '../components/DataTable';
import { Plus, Edit2, Trash2, X } from 'lucide-react';
import type { ColumnDef } from '@tanstack/react-table';

interface Server {
  id: number;
  name: string;
  ip: string;
  os: string;
  location: string;
  CreatedAt: string;
}

export const ServersOverview: React.FC = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingServer, setEditingServer] = useState<Server | null>(null);

  // Form states
  const [name, setName] = useState('');
  const [ip, setIp] = useState('');
  const [os, setOs] = useState('');
  const [location, setLocation] = useState('');
  const [error, setError] = useState<string | null>(null);

  // Fetch servers
  const { data: servers = [], isLoading } = useQuery<Server[]>({
    queryKey: ['servers'],
    queryFn: async () => {
      const response = await api.get('/api/servers/');
      return response.data;
    },
  });

  // Mutations
  const createMutation = useMutation({
    mutationFn: async (newServer: Omit<Server, 'id' | 'CreatedAt'>) => {
      const response = await api.post('/api/servers/', newServer);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
      closeModal();
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to create server');
    }
  });

  const updateMutation = useMutation({
    mutationFn: async (updatedServer: Server) => {
      const response = await api.put(`/api/servers/${updatedServer.id}`, {
        name: updatedServer.name,
        ip: updatedServer.ip,
        os: updatedServer.os,
        location: updatedServer.location,
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
      closeModal();
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to update server');
    }
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/api/servers/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
    },
    onError: (err: any) => {
      alert(err.response?.data?.error || 'Failed to delete server');
    }
  });

  const openCreateModal = () => {
    setEditingServer(null);
    setName('');
    setIp('');
    setOs('');
    setLocation('');
    setError(null);
    setIsModalOpen(true);
  };

  const openEditModal = (server: Server) => {
    setEditingServer(server);
    setName(server.name);
    setIp(server.ip);
    setOs(server.os);
    setLocation(server.location);
    setError(null);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setEditingServer(null);
    setName('');
    setIp('');
    setOs('');
    setLocation('');
    setError(null);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !ip.trim() || !os.trim() || !location.trim()) {
      setError('All fields are required');
      return;
    }

    const payload = {
      name,
      ip,
      os,
      location,
    };

    if (editingServer) {
      updateMutation.mutate({
        ...editingServer,
        ...payload,
      });
    } else {
      createMutation.mutate(payload);
    }
  };

  const handleDelete = (id: number, serverName: string) => {
    if (confirm(`Are you sure you want to delete server "${serverName}"? Systems currently pointing to this server ID will reference an unavailable server.`)) {
      deleteMutation.mutate(id);
    }
  };

  const isViewer = user?.role === 'Viewer';

  // Table Columns Definition
  const columns = React.useMemo<ColumnDef<Server>[]>(() => [
    {
      accessorKey: 'name',
      header: 'Server Name',
      cell: info => <span style={{ fontWeight: 600 }}>{info.getValue() as string}</span>,
    },
    {
      accessorKey: 'ip',
      header: 'IP Address',
      cell: info => <code style={{ background: 'var(--bg-tertiary)', padding: '2px 6px', borderRadius: '4px' }}>{info.getValue() as string}</code>,
    },
    {
      accessorKey: 'os',
      header: 'Operating System',
      cell: info => <span>{info.getValue() as string}</span>,
    },
    {
      accessorKey: 'location',
      header: 'Data Center Location',
      cell: info => <span>{info.getValue() as string}</span>,
    },
    {
      id: 'actions',
      header: 'Actions',
      cell: ({ row }) => {
        const server = row.original;
        return (
          <div style={{ display: 'flex', gap: '8px' }}>
            {!isViewer && (
              <>
                <button
                  onClick={() => openEditModal(server)}
                  className="btn btn-secondary btn-sm"
                  style={{ padding: '6px' }}
                >
                  <Edit2 size={14} />
                </button>
                <button
                  onClick={() => handleDelete(server.id, server.name)}
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
  ], [isViewer]);

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: '64px' }}>
        <div style={{ width: '45px', height: '45px', border: '3px solid var(--border)', borderTopColor: 'var(--accent)', borderRadius: '50%', animation: 'spin 1s linear infinite', margin: '0 auto 16px' }} />
        <p style={{ color: 'var(--text-secondary)' }}>Loading servers list...</p>
      </div>
    );
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div>
          <p style={{ color: 'var(--text-secondary)' }}>Manage the physical/virtual host servers that systems deploy on.</p>
        </div>
        {!isViewer && (
          <button onClick={openCreateModal} className="btn btn-primary">
            <Plus size={16} />
            <span>Add Server</span>
          </button>
        )}
      </div>

      <div className="card" style={{ padding: '20px' }}>
        <DataTable columns={columns} data={servers} searchPlaceholder="Search servers by name, IP, location..." />
      </div>

      {/* Create / Edit Modal */}
      {isModalOpen && (
        <div className="modal-overlay">
          <div className="modal-content">
            <div className="modal-header">
              <h3 className="modal-title">{editingServer ? 'Edit Server' : 'Add New Server'}</h3>
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
                  <label className="form-label">Server Name</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. K8S-Worker-Prod-01"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">IP Address</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. 10.120.45.12"
                    value={ip}
                    onChange={(e) => setIp(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Operating System</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. Ubuntu 22.04 LTS, RHEL 9.1"
                    value={os}
                    onChange={(e) => setOs(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Location / Data Center</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. AWS ap-southeast-1, Jakarta-DC2"
                    value={location}
                    onChange={(e) => setLocation(e.target.value)}
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={closeModal}>Cancel</button>
                <button type="submit" className="btn btn-primary">
                  {editingServer ? 'Save Changes' : 'Create Server'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
