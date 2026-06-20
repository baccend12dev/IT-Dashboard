import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import { api } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { FileText, Save } from 'lucide-react';

interface System {
  id: number;
  name: string;
}

export const GlobalAddDocumentation: React.FC = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isViewer = user?.role === 'Viewer';

  // Form states
  const [selectedSystemId, setSelectedSystemId] = useState('');
  const [title, setTitle] = useState('');
  const [category, setCategory] = useState('Technical Flow');
  const [content, setContent] = useState('');
  const [error, setError] = useState<string | null>(null);

  // Fetch systems list for selection dropdown
  const { data: systems = [], isLoading: isLoadingSystems } = useQuery<System[]>({
    queryKey: ['systems'],
    queryFn: async () => {
      const response = await api.get('/api/systems/');
      return response.data;
    },
  });

  const createDocMutation = useMutation({
    mutationFn: async (payload: { systemId: number; title: string; category: string; content: string }) => {
      const response = await api.post(`/api/systems/${payload.systemId}/documentations`, {
        title: payload.title,
        category: payload.category,
        content: payload.content,
      });
      return response.data;
    },
    onSuccess: (_, variables) => {
      // Invalidate target system documentations
      queryClient.invalidateQueries({ queryKey: ['docs', String(variables.systemId)] });
      // Redirect to system details docs tab
      navigate({
        to: `/dashboard/systems/${variables.systemId}`,
      });
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to save documentation writeup');
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (isViewer) {
      setError('Viewer accounts are restricted from creating documentation');
      return;
    }

    if (!selectedSystemId) {
      setError('Please select an IT System');
      return;
    }
    if (!title.trim()) {
      setError('Document title is required');
      return;
    }
    if (!content.trim()) {
      setError('Document content cannot be empty');
      return;
    }

    createDocMutation.mutate({
      systemId: Number(selectedSystemId),
      title,
      category,
      content,
    });
  };

  const handleCancel = () => {
    navigate({ to: '/dashboard' });
  };

  if (isLoadingSystems) {
    return (
      <div style={{ textAlign: 'center', padding: '64px' }}>
        <div style={{ width: '45px', height: '45px', border: '3px solid var(--border)', borderTopColor: 'var(--accent)', borderRadius: '50%', animation: 'spin 1s linear infinite', margin: '0 auto 16px' }} />
        <p style={{ color: 'var(--text-secondary)' }}>Loading systems database...</p>
      </div>
    );
  }

  return (
    <div className="card" style={{ maxWidth: '800px', margin: '0 auto' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', borderBottom: '1px solid var(--border)', paddingBottom: '16px', marginBottom: '24px' }}>
        <div style={{ width: '36px', height: '36px', background: 'var(--accent-glow)', color: 'var(--accent-light)', borderRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <FileText size={18} />
        </div>
        <div>
          <h3 style={{ fontSize: '18px' }}>Create Global System Documentation</h3>
          <p style={{ color: 'var(--text-secondary)', fontSize: '13px', marginTop: '2px' }}>Publish a new technical writeup or flow chart guide to any cataloged system.</p>
        </div>
      </div>

      {isViewer && (
        <div className="login-error" style={{ marginBottom: '20px' }}>
          <span>Warning: You are logged in as a Viewer. You cannot publish documentation.</span>
        </div>
      )}

      {error && (
        <div className="login-error" style={{ marginBottom: '20px' }}>
          <span>{error}</span>
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
          <div className="form-group">
            <label className="form-label">Select IT System *</label>
            <select
              className="form-select"
              value={selectedSystemId}
              onChange={(e) => setSelectedSystemId(e.target.value)}
              disabled={isViewer || systems.length === 0}
            >
              <option value="" disabled>-- Pick a system --</option>
              {systems.map(sys => (
                <option key={sys.id} value={sys.id}>{sys.name}</option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label className="form-label">Document Category</label>
            <select
              className="form-select"
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              disabled={isViewer}
            >
              <option value="Business Flow">Business Flow</option>
              <option value="Technical Flow">Technical Flow</option>
              <option value="API Documentation">API Documentation</option>
              <option value="Database Documentation">Database Documentation</option>
              <option value="Deployment Guide">Deployment Guide</option>
              <option value="User Manual">User Manual</option>
            </select>
          </div>
        </div>

        <div className="form-group">
          <label className="form-label">Document Title *</label>
          <input
            type="text"
            className="form-input"
            placeholder="e.g. RabbitMQ Event Driven Architecture and Consumer Setup"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={isViewer}
          />
        </div>

        <div className="form-group">
          <label className="form-label">Markdown Content / Text Writeup *</label>
          <textarea
            className="form-textarea"
            style={{ minHeight: '300px', fontFamily: 'monospace' }}
            placeholder="Write detailed documentation here. Support markdown or plain text headers."
            value={content}
            onChange={(e) => setContent(e.target.value)}
            disabled={isViewer}
          />
        </div>

        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px', borderTop: '1px solid var(--border)', paddingTop: '20px', marginTop: '20px' }}>
          <button type="button" className="btn btn-secondary" onClick={handleCancel}>Cancel</button>
          <button
            type="submit"
            className="btn btn-primary"
            disabled={isViewer || createDocMutation.isPending || systems.length === 0}
          >
            <Save size={16} />
            <span>{createDocMutation.isPending ? 'Publishing...' : 'Publish Document'}</span>
          </button>
        </div>
      </form>
    </div>
  );
};
