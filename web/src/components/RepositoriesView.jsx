import { motion } from 'framer-motion';
import { CheckCircle, XCircle, Clock, FolderGit2 } from 'lucide-react';

export const RepositoriesView = ({ repositories }) => {
  const getStatusIcon = (status) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="w-5 h-5 text-green-400" />;
      case 'failed':
        return <XCircle className="w-5 h-5 text-red-400" />;
      default:
        return <Clock className="w-5 h-5 text-gray-400" />;
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'success':
        return 'border-green-500/30 bg-green-500/5';
      case 'failed':
        return 'border-red-500/30 bg-red-500/5';
      default:
        return 'border-gray-500/30 bg-gray-500/5';
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-text-primary mb-2">Repositories</h2>
        <p className="text-text-tertiary">
          Cloned repositories in your workspace
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {repositories.map((repo, index) => (
          <motion.div
            key={repo.name}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3, delay: index * 0.1 }}
            className={`card border-2 ${getStatusColor(repo.status)}`}
          >
            {/* Header */}
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-surface-hover rounded-lg flex items-center justify-center">
                  <FolderGit2 className="w-5 h-5 text-primary-400" />
                </div>
                <div>
                  <h3 className="font-semibold text-text-primary">{repo.name}</h3>
                  <p className="text-xs text-text-tertiary mt-0.5">
                    {repo.duration_seconds ? `${repo.duration_seconds.toFixed(1)}s` : '-'}
                  </p>
                </div>
              </div>
              {getStatusIcon(repo.status)}
            </div>

            {/* Path */}
            <div className="mb-3">
              <p className="text-xs text-text-tertiary mb-1">Path</p>
              <p className="text-sm text-text-secondary font-mono break-all bg-surface-base px-2 py-1 rounded">
                {repo.path || 'N/A'}
              </p>
            </div>

            {/* URL */}
            <div className="mb-3">
              <p className="text-xs text-text-tertiary mb-1">Repository URL</p>
              <p className="text-sm text-text-secondary font-mono break-all bg-surface-base px-2 py-1 rounded">
                {repo.url || 'N/A'}
              </p>
            </div>

            {/* Setup Commands */}
            {repo.setup_commands && repo.setup_commands.length > 0 && (
              <div>
                <p className="text-xs text-text-tertiary mb-2">Setup Commands</p>
                <div className="space-y-1">
                  {repo.setup_commands.map((cmd, idx) => (
                    <div
                      key={idx}
                      className="text-xs text-text-secondary font-mono bg-surface-base px-2 py-1 rounded"
                    >
                      {cmd}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Error Message */}
            {repo.error && (
              <div className="mt-3 p-2 bg-red-500/10 border border-red-500/20 rounded text-xs text-red-400">
                {repo.error}
              </div>
            )}

            {/* Status Badge */}
            <div className="mt-4 pt-3 border-t border-surface-border">
              <span className={`badge ${
                repo.status === 'success' ? 'badge-success' :
                repo.status === 'failed' ? 'badge-error' :
                'badge-info'
              }`}>
                {repo.status || 'unknown'}
              </span>
            </div>
          </motion.div>
        ))}
      </div>

      {repositories.length === 0 && (
        <div className="text-center py-16 text-text-tertiary">
          <FolderGit2 className="w-12 h-12 mx-auto mb-4 opacity-50" />
          <p>No repositories initialized yet</p>
        </div>
      )}
    </div>
  );
};
