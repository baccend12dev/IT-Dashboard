import React, { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams, Link } from '@tanstack/react-router';
import { api } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { ChevronLeft, Plus, Edit2, Trash2, X, Terminal, Server as ServerIcon, Clock, Award, FileText } from 'lucide-react';

// Interfaces matching backend models
interface Server {
  id: number;
  name: string;
  ip: string;
  os: string;
  location: string;
}

interface System {
  id: number;
  name: string;
  type: string;
  links: string;
  server_id: number;
  status: string;
  description: string;
  Server?: Server;
}

interface Note {
  id: number;
  system_id: number;
  title: string;
  content: string;
  created_at: string;
}

interface FeatureRequest {
  id: number;
  system_id: number;
  title: string;
  description: string;
  status: string; // Pending, Approved, In Progress, Completed, Rejected
  created_at: string;
}

interface Documentation {
  id: number;
  system_id: number;
  title: string;
  category: string; // Business Flow, Technical Flow, API Documentation, Database Documentation, Deployment Guide, User Manual
  content: string;
  created_at: string;
}

export const SystemDetail: React.FC = () => {
  const { systemId } = useParams({ strict: false });
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const isViewer = user?.role === 'Viewer';

  // Active tab state
  const [activeTab, setActiveTab] = useState<'info' | 'docs' | 'notes' | 'features'>('info');

  // Sub-resource Modal states
  const [activeModal, setActiveModal] = useState<'none' | 'note' | 'feature' | 'doc'>('none');
  const [editingItem, setEditingItem] = useState<any>(null); // holds Note/FeatureRequest/Documentation model

  // Shared form inputs
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [category, setCategory] = useState('Technical Flow'); // doc category default
  const [status, setStatus] = useState('Pending'); // feature status default
  const [error, setError] = useState<string | null>(null);

  // Selected document for preview
  const [selectedDoc, setSelectedDoc] = useState<Documentation | null>(null);

  // Queries
  const { data: system } = useQuery<System>({
    queryKey: ['system', systemId],
    queryFn: async () => {
      const response = await api.get(`/api/systems/${systemId}`);
      return response.data;
    },
  });

  const { data: notes } = useQuery<Note[]>({
    queryKey: ['notes', systemId],
    queryFn: async () => {
      const response = await api.get(`/api/systems/${systemId}/notes`);
      return response.data;
    },
    enabled: !!systemId,
  });
  const noteList = notes || [];

  const { data: features } = useQuery<FeatureRequest[]>({
    queryKey: ['features', systemId],
    queryFn: async () => {
      const response = await api.get(`/api/systems/${systemId}/feature-requests`);
      return response.data;
    },
    enabled: !!systemId,
  });
  const featureList = features || [];

  const { data: docs } = useQuery<Documentation[]>({
    queryKey: ['docs', systemId],
    queryFn: async () => {
      const response = await api.get(`/api/systems/${systemId}/documentations`);
      return response.data;
    },
    enabled: !!systemId,
  });
  const docList = docs || [];

  // Dynamic document auto-selection fallback
  useEffect(() => {
    if (docList.length > 0 && !selectedDoc) {
      setSelectedDoc(docList[0]);
    }
  }, [docList, selectedDoc]);

  // Mutations
  const noteMutation = useMutation({
    mutationFn: async (payload: { id?: number; title: string; content: string }) => {
      if (payload.id) {
        const response = await api.put(`/api/notes/${payload.id}`, payload);
        return response.data;
      } else {
        const response = await api.post(`/api/systems/${systemId}/notes`, payload);
        return response.data;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notes', systemId] });
      closeModal();
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to save note');
    }
  });

  const featureMutation = useMutation({
    mutationFn: async (payload: { id?: number; title: string; description: string; status?: string }) => {
      if (payload.id) {
        const response = await api.put(`/api/feature-requests/${payload.id}`, {
          title: payload.title,
          description: payload.description,
          status: payload.status,
        });
        return response.data;
      } else {
        const response = await api.post(`/api/systems/${systemId}/feature-requests`, {
          title: payload.title,
          description: payload.description,
        });
        return response.data;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['features', systemId] });
      closeModal();
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to save feature request');
    }
  });

  const docMutation = useMutation({
    mutationFn: async (payload: { id?: number; title: string; category: string; content: string }) => {
      if (payload.id) {
        const response = await api.put(`/api/documentations/${payload.id}`, payload);
        return response.data;
      } else {
        const response = await api.post(`/api/systems/${systemId}/documentations`, payload);
        return response.data;
      }
    },
    onSuccess: (data: Documentation) => {
      queryClient.invalidateQueries({ queryKey: ['docs', systemId] });
      setSelectedDoc(data);
      closeModal();
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to save documentation');
    }
  });

  const deleteItemMutation = useMutation({
    mutationFn: async (payload: { type: 'note' | 'feature' | 'doc'; id: number }) => {
      let endpoint = '';
      if (payload.type === 'note') endpoint = `/api/notes/${payload.id}`;
      else if (payload.type === 'feature') endpoint = `/api/feature-requests/${payload.id}`;
      else if (payload.type === 'doc') endpoint = `/api/documentations/${payload.id}`;
      
      await api.delete(endpoint);
    },
    onSuccess: (_, variables) => {
      if (variables.type === 'note') {
        queryClient.invalidateQueries({ queryKey: ['notes', systemId] });
      } else if (variables.type === 'feature') {
        queryClient.invalidateQueries({ queryKey: ['features', systemId] });
      } else if (variables.type === 'doc') {
        queryClient.invalidateQueries({ queryKey: ['docs', systemId] });
        if (selectedDoc && selectedDoc.id === variables.id) {
          setSelectedDoc(null);
        }
      }
    },
    onError: (err: any) => {
      alert(err.response?.data?.error || 'Failed to delete record');
    }
  });

  // Modal open helpers
  const openNoteModal = (note?: Note) => {
    if (note) {
      setEditingItem(note);
      setTitle(note.title);
      setContent(note.content);
    } else {
      setEditingItem(null);
      setTitle('');
      setContent('');
    }
    setError(null);
    setActiveModal('note');
  };

  const openFeatureModal = (feature?: FeatureRequest) => {
    if (feature) {
      setEditingItem(feature);
      setTitle(feature.title);
      setContent(feature.description);
      setStatus(feature.status);
    } else {
      setEditingItem(null);
      setTitle('');
      setContent('');
      setStatus('Pending');
    }
    setError(null);
    setActiveModal('feature');
  };

  const openDocModal = (doc?: Documentation) => {
    if (doc) {
      setEditingItem(doc);
      setTitle(doc.title);
      setCategory(doc.category);
      setContent(doc.content);
    } else {
      setEditingItem(null);
      setTitle('');
      setCategory('Technical Flow');
      setContent('');
    }
    setError(null);
    setActiveModal('doc');
  };

  const closeModal = () => {
    setActiveModal('none');
    setEditingItem(null);
    setTitle('');
    setContent('');
    setCategory('Technical Flow');
    setStatus('Pending');
    setError(null);
  };

  // Submit handlers
  const handleNoteSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) {
      setError('Title and Content are required');
      return;
    }
    noteMutation.mutate({ id: editingItem?.id, title, content });
  };

  const handleFeatureSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) {
      setError('Title and Description are required');
      return;
    }
    featureMutation.mutate({ id: editingItem?.id, title, description: content, status });
  };

  const handleDocSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !category.trim() || !content.trim()) {
      setError('All fields are required');
      return;
    }
    docMutation.mutate({ id: editingItem?.id, title, category, content });
  };

  const handleDeleteItem = (type: 'note' | 'feature' | 'doc', id: number, itemTitle: string) => {
    if (confirm(`Are you sure you want to delete this ${type} "${itemTitle}"?`)) {
      deleteItemMutation.mutate({ type, id });
    }
  };

  if (!system) {
    return (
      <div className="card" style={{ textAlign: 'center', padding: '48px' }}>
        <h3>System not found</h3>
        <p style={{ color: 'var(--text-muted)', margin: '12px 0 24px' }}>The system record you requested does not exist or has been deleted.</p>
        <Link to="/dashboard" className="btn btn-secondary">
          <ChevronLeft size={16} />
          Back to Systems
        </Link>
      </div>
    );
  }

  return (
    <div>
      {/* Header breadcrumb */}
      <div style={{ marginBottom: '24px' }}>
        <Link to="/dashboard" style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', fontSize: '14px', color: 'var(--text-secondary)' }}>
          <ChevronLeft size={16} />
          <span>Back to Systems List</span>
        </Link>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginTop: '12px' }}>
          <div>
            <h2 style={{ fontSize: '28px', display: 'flex', alignItems: 'center', gap: '12px' }}>
              {system.name}
              <span className={`badge ${system.status === 'Active' || system.status === 'Online' ? 'badge-success' : 'badge-warning'}`} style={{ fontSize: '12px' }}>
                {system.status}
              </span>
            </h2>
            <p style={{ color: 'var(--text-secondary)', marginTop: '4px' }}>Deployment Type: <code>{system.type}</code></p>
          </div>
          <a href={system.links} target="_blank" rel="noreferrer" className="btn btn-secondary">
            Access Application Endpoint
          </a>
        </div>
      </div>

      {/* Tabs list */}
      <div className="tabs-header">
        <button className={`tab-btn ${activeTab === 'info' ? 'active' : ''}`} onClick={() => setActiveTab('info')}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <ServerIcon size={16} />
            <span>Specifications & Servers</span>
          </div>
        </button>
        <button className={`tab-btn ${activeTab === 'docs' ? 'active' : ''}`} onClick={() => setActiveTab('docs')}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <FileText size={16} />
            <span>System Documentation ({docList.length})</span>
          </div>
        </button>
        <button className={`tab-btn ${activeTab === 'notes' ? 'active' : ''}`} onClick={() => setActiveTab('notes')}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Clock size={16} />
            <span>Developer Notes ({noteList.length})</span>
          </div>
        </button>
        <button className={`tab-btn ${activeTab === 'features' ? 'active' : ''}`} onClick={() => setActiveTab('features')}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Award size={16} />
            <span>Feature Backlog ({featureList.length})</span>
          </div>
        </button>
      </div>

      {/* Tab Contents */}
      <div style={{ minHeight: '400px' }}>
        {activeTab === 'info' && (
          <div className="grid-2">
            <div className="card">
              <h3 style={{ marginBottom: '20px', display: 'flex', alignItems: 'center', gap: '8px', borderBottom: '1px solid var(--border)', paddingBottom: '10px' }}>
                <Terminal size={18} style={{ color: 'var(--accent)' }} />
                System Specification
              </h3>
              <div className="detail-item">
                <span className="detail-label">Name</span>
                <span className="detail-value">{system.name}</span>
              </div>
              <div className="detail-item">
                <span className="detail-label">Stack Type</span>
                <span className="detail-value"><code>{system.type}</code></span>
              </div>
              <div className="detail-item">
                <span className="detail-label">Application Endpoint</span>
                <span className="detail-value"><a href={system.links} target="_blank" rel="noreferrer">{system.links}</a></span>
              </div>
              <div className="detail-item">
                <span className="detail-label">Description</span>
                <span className="detail-value" style={{ color: 'var(--text-secondary)' }}>{system.description}</span>
              </div>
            </div>

            <div className="card">
              <h3 style={{ marginBottom: '20px', display: 'flex', alignItems: 'center', gap: '8px', borderBottom: '1px solid var(--border)', paddingBottom: '10px' }}>
                <ServerIcon size={18} style={{ color: 'var(--success)' }} />
                Host Node Specification
              </h3>
              {system.Server ? (
                <div>
                  <div className="detail-item">
                    <span className="detail-label">Server Name</span>
                    <span className="detail-value">{system.Server.name}</span>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">IP Address</span>
                    <span className="detail-value"><code>{system.Server.ip}</code></span>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">Operating System</span>
                    <span className="detail-value">{system.Server.os}</span>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">Physical Datacenter Location</span>
                    <span className="detail-value">{system.Server.location}</span>
                  </div>
                </div>
              ) : (
                <div style={{ padding: '24px 0', textAlign: 'center', color: 'var(--text-muted)' }}>
                  No host server details loaded. Linked server ID: {system.server_id}
                </div>
              )}
            </div>
          </div>
        )}

        {activeTab === 'docs' && (
          <div style={{ display: 'grid', gridTemplateColumns: '280px 1fr', gap: '24px', alignItems: 'start' }}>
            {/* Sidebar for Docs List */}
            <div className="card" style={{ padding: '16px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
                <h4 style={{ fontSize: '14px', textTransform: 'uppercase', color: 'var(--text-muted)' }}>Writeups</h4>
                {!isViewer && (
                  <button onClick={() => openDocModal()} className="btn btn-secondary btn-sm" style={{ padding: '4px 8px', fontSize: '11px' }}>
                    <Plus size={12} /> Add
                  </button>
                )}
              </div>
              
              {docList.length > 0 ? (
                docList.map(doc => (
                  <button
                    key={doc.id}
                    onClick={() => setSelectedDoc(doc)}
                    style={{
                      display: 'flex',
                      flexDirection: 'column',
                      alignItems: 'start',
                      width: '100%',
                      padding: '10px 12px',
                      background: selectedDoc?.id === doc.id ? 'var(--accent-glow)' : 'transparent',
                      border: '1px solid',
                      borderColor: selectedDoc?.id === doc.id ? 'rgba(99, 102, 241, 0.2)' : 'transparent',
                      borderRadius: 'var(--radius-md)',
                      color: selectedDoc?.id === doc.id ? 'var(--accent-light)' : 'var(--text-secondary)',
                      textAlign: 'left',
                      cursor: 'pointer',
                      transition: 'var(--transition-smooth)'
                    }}
                  >
                    <span style={{ fontWeight: 600, fontSize: '13px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', width: '100%' }}>{doc.title}</span>
                    <span style={{ fontSize: '10px', opacity: 0.7, marginTop: '2px', background: 'var(--bg-tertiary)', padding: '1px 4px', borderRadius: '3px' }}>{doc.category}</span>
                  </button>
                ))
              ) : (
                <div style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '16px 0', fontSize: '13px' }}>
                  No documents found.
                </div>
              )}
            </div>

            {/* Document Content View */}
            <div className="card">
              {selectedDoc ? (
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', borderBottom: '1px solid var(--border)', paddingBottom: '16px', marginBottom: '16px' }}>
                    <div>
                      <h3 style={{ fontSize: '20px' }}>{selectedDoc.title}</h3>
                      <span className="badge badge-info" style={{ marginTop: '6px', fontSize: '10px' }}>{selectedDoc.category}</span>
                    </div>
                    {!isViewer && (
                      <div style={{ display: 'flex', gap: '8px' }}>
                        <button onClick={() => openDocModal(selectedDoc)} className="btn btn-secondary btn-sm">
                          <Edit2 size={12} />
                          <span>Edit</span>
                        </button>
                        <button onClick={() => handleDeleteItem('doc', selectedDoc.id, selectedDoc.title)} className="btn btn-danger btn-sm">
                          <Trash2 size={12} />
                          <span>Delete</span>
                        </button>
                      </div>
                    )}
                  </div>
                  <div className="markdown-preview">{selectedDoc.content}</div>
                </div>
              ) : (
                <div style={{ textAlign: 'center', padding: '64px 0', color: 'var(--text-muted)' }}>
                  Select a document writeup from the list on the left to read its contents.
                </div>
              )}
            </div>
          </div>
        )}

        {activeTab === 'notes' && (
          <div className="card">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
              <h3 style={{ fontSize: '16px' }}>System Logs and Timeline Update</h3>
              {!isViewer && (
                <button onClick={() => openNoteModal()} className="btn btn-primary btn-sm">
                  <Plus size={14} />
                  <span>Add Note</span>
                </button>
              )}
            </div>

            {noteList.length > 0 ? (
              <div className="timeline">
                {noteList.map(note => (
                  <div className="timeline-item" key={note.id}>
                    <div className="timeline-dot" />
                    <div className="timeline-content">
                      <div className="timeline-header">
                        <h4 style={{ fontSize: '14px', fontWeight: 600 }}>{note.title}</h4>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                          <span className="timeline-meta">{new Date(note.created_at).toLocaleString()}</span>
                          {!isViewer && (
                            <div style={{ display: 'flex', gap: '4px' }}>
                              <button onClick={() => openNoteModal(note)} className="btn btn-secondary btn-sm" style={{ padding: '3px' }}>
                                <Edit2 size={10} />
                              </button>
                              <button onClick={() => handleDeleteItem('note', note.id, note.title)} className="btn btn-danger btn-sm" style={{ padding: '3px' }}>
                                <Trash2 size={10} />
                              </button>
                            </div>
                          )}
                        </div>
                      </div>
                      <p className="timeline-body" style={{ color: 'var(--text-secondary)', whiteSpace: 'pre-wrap' }}>{note.content}</p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div style={{ textAlign: 'center', padding: '48px 0', color: 'var(--text-muted)' }}>
                No developer notes registered for this system yet.
              </div>
            )}
          </div>
        )}

        {activeTab === 'features' && (
          <div className="card">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
              <h3 style={{ fontSize: '16px' }}>System Feature Backlog</h3>
              {!isViewer && (
                <button onClick={() => openFeatureModal()} className="btn btn-primary btn-sm">
                  <Plus size={14} />
                  <span>Request Feature</span>
                </button>
              )}
            </div>

            {featureList.length > 0 ? (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                {featureList.map(feat => {
                  let statusColor = 'badge-info';
                  if (feat.status === 'Completed' || feat.status === 'Approved') statusColor = 'badge-success';
                  if (feat.status === 'In Progress') statusColor = 'badge-warning';
                  if (feat.status === 'Rejected') statusColor = 'badge-danger';
                  
                  return (
                    <div key={feat.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', padding: '16px', background: 'var(--bg-tertiary)', borderRadius: 'var(--radius-md)', border: '1px solid var(--border)' }}>
                      <div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                          <h4 style={{ fontSize: '15px', fontWeight: 600 }}>{feat.title}</h4>
                          <span className={`badge ${statusColor}`} style={{ fontSize: '10px' }}>{feat.status}</span>
                        </div>
                        <p style={{ color: 'var(--text-secondary)', fontSize: '13px', marginTop: '6px', whiteSpace: 'pre-wrap' }}>{feat.description}</p>
                        <span style={{ fontSize: '11px', color: 'var(--text-muted)', display: 'block', marginTop: '10px' }}>
                          Created: {new Date(feat.created_at).toLocaleDateString()}
                        </span>
                      </div>
                      {!isViewer && (
                        <div style={{ display: 'flex', gap: '6px' }}>
                          <button onClick={() => openFeatureModal(feat)} className="btn btn-secondary btn-sm" style={{ padding: '6px' }}>
                            <Edit2 size={12} />
                          </button>
                          <button onClick={() => handleDeleteItem('feature', feat.id, feat.title)} className="btn btn-danger btn-sm" style={{ padding: '6px' }}>
                            <Trash2 size={12} />
                          </button>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            ) : (
              <div style={{ textAlign: 'center', padding: '48px 0', color: 'var(--text-muted)' }}>
                No feature requests mapped.
              </div>
            )}
          </div>
        )}
      </div>

      {/* CRUD Modals */}
      {activeModal === 'note' && (
        <div className="modal-overlay">
          <div className="modal-content">
            <div className="modal-header">
              <h3 className="modal-title">{editingItem ? 'Edit Developer Note' : 'Add Developer Note'}</h3>
              <button className="modal-close" onClick={closeModal}><X size={18} /></button>
            </div>
            <form onSubmit={handleNoteSubmit}>
              <div className="modal-body">
                {error && <div className="login-error" style={{ marginBottom: '16px' }}>{error}</div>}
                <div className="form-group">
                  <label className="form-label">Note Title</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. Completed migration to v2 databases"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Note Details / Content</label>
                  <textarea
                    className="form-textarea"
                    placeholder="Provide details about updates, patch updates, or database fixes."
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={closeModal}>Cancel</button>
                <button type="submit" className="btn btn-primary">Save Note</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {activeModal === 'feature' && (
        <div className="modal-overlay">
          <div className="modal-content">
            <div className="modal-header">
              <h3 className="modal-title">{editingItem ? 'Edit Feature Request' : 'Request New Feature'}</h3>
              <button className="modal-close" onClick={closeModal}><X size={18} /></button>
            </div>
            <form onSubmit={handleFeatureSubmit}>
              <div className="modal-body">
                {error && <div className="login-error" style={{ marginBottom: '16px' }}>{error}</div>}
                <div className="form-group">
                  <label className="form-label">Feature Title</label>
                  <input
                    type="text"
                    className="form-input"
                    placeholder="e.g. Excel reports exporter integration"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                  />
                </div>
                {editingItem && (
                  <div className="form-group">
                    <label className="form-label">Status</label>
                    <select
                      className="form-select"
                      value={status}
                      onChange={(e) => setStatus(e.target.value)}
                    >
                      <option value="Pending">Pending</option>
                      <option value="Approved">Approved</option>
                      <option value="In Progress">In Progress</option>
                      <option value="Completed">Completed</option>
                      <option value="Rejected">Rejected</option>
                    </select>
                  </div>
                )}
                <div className="form-group">
                  <label className="form-label">Feature Description</label>
                  <textarea
                    className="form-textarea"
                    placeholder="Describe what requirements, triggers, and mock endpoints this feature requires."
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={closeModal}>Cancel</button>
                <button type="submit" className="btn btn-primary">Save Feature</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {activeModal === 'doc' && (
        <div className="modal-overlay">
          <div className="modal-content large">
            <div className="modal-header">
              <h3 className="modal-title">{editingItem ? 'Edit Document Writeup' : 'Add System Document Writeup'}</h3>
              <button className="modal-close" onClick={closeModal}><X size={18} /></button>
            </div>
            <form onSubmit={handleDocSubmit}>
              <div className="modal-body">
                {error && <div className="login-error" style={{ marginBottom: '16px' }}>{error}</div>}
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 200px', gap: '16px' }}>
                  <div className="form-group">
                    <label className="form-label">Document Title</label>
                    <input
                      type="text"
                      className="form-input"
                      placeholder="e.g. Deployment setup instructions"
                      value={title}
                      onChange={(e) => setTitle(e.target.value)}
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">Category</label>
                    <select
                      className="form-select"
                      value={category}
                      onChange={(e) => setCategory(e.target.value)}
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
                  <label className="form-label">Markdown / Text Writeup Content</label>
                  <textarea
                    className="form-textarea"
                    style={{ minHeight: '220px', fontFamily: 'monospace' }}
                    placeholder="Write detailed documentation here. Support plain text / markdown formats."
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={closeModal}>Cancel</button>
                <button type="submit" className="btn btn-primary">Save Document</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
