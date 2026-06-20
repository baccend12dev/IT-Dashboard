import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from '@tanstack/react-router';
import { api } from '../services/api';
import { DataTable } from '../components/DataTable';
import { Eye } from 'lucide-react';
import type { ColumnDef } from '@tanstack/react-table';

interface System {
  id: number;
  name: string;
}

interface FeatureRequest {
  id: number;
  system_id: number;
  system?: System;
  title: string;
  description: string;
  status: string;
  created_at: string;
}

export const PendingRequests: React.FC = () => {
  const { data: requests = [], isLoading } = useQuery<FeatureRequest[]>({
    queryKey: ['pending-features'],
    queryFn: async () => {
      const response = await api.get('/api/feature-requests/pending');
      return response.data;
    },
  });

  const columns = React.useMemo<ColumnDef<FeatureRequest>[]>(() => [
    {
      id: 'system_name',
      header: 'System Name',
      cell: ({ row }) => {
        const item = row.original;
        if (item.system) {
          return (
            <Link
              to="/dashboard/systems/$systemId"
              params={{ systemId: String(item.system.id) }}
              style={{ fontWeight: 600, color: 'var(--accent)' }}
            >
              {item.system.name}
            </Link>
          );
        }
        return <span style={{ color: 'var(--text-secondary)' }}>System ID: {item.system_id}</span>;
      },
    },
    {
      accessorKey: 'title',
      header: 'Feature Title',
      cell: info => <span style={{ fontWeight: 500 }}>{info.getValue() as string}</span>,
    },
    {
      accessorKey: 'description',
      header: 'Description',
      cell: info => <span style={{ color: 'var(--text-secondary)' }}>{info.getValue() as string}</span>,
    },
    {
      accessorKey: 'created_at',
      header: 'Requested Date',
      cell: info => {
        const val = info.getValue() as string;
        if (!val) return '-';
        return <span>{new Date(val).toLocaleDateString('id-ID', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit'
        })}</span>;
      },
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: info => {
        const val = info.getValue() as string;
        return <span className="badge badge-warning">{val}</span>;
      },
    },
    {
      id: 'actions',
      header: 'Actions',
      cell: ({ row }) => {
        const item = row.original;
        return (
          <div style={{ display: 'flex', gap: '8px' }}>
            <Link
              to="/dashboard/systems/$systemId"
              params={{ systemId: String(item.system_id) }}
              className="btn btn-secondary btn-sm"
              style={{ padding: '6px' }}
              title="View System Details"
            >
              <Eye size={14} />
            </Link>
          </div>
        );
      },
    },
  ], []);

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: '64px' }}>
        <div style={{ width: '45px', height: '45px', border: '3px solid var(--border)', borderTopColor: 'var(--accent)', borderRadius: '50%', animation: 'spin 1s linear infinite', margin: '0 auto 16px' }} />
        <p style={{ color: 'var(--text-secondary)' }}>Loading pending feature requests...</p>
      </div>
    );
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div>
          <p style={{ color: 'var(--text-secondary)' }}>
            List of all feature requests that are currently pending implementation across all enterprise systems.
          </p>
        </div>
      </div>

      <div className="card" style={{ padding: '20px' }}>
        <DataTable columns={columns} data={requests} searchPlaceholder="Search pending features..." />
      </div>
    </div>
  );
};
