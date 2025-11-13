import { useState } from 'react';
import { Upload, FileText, CheckCircle, AlertCircle, Loader2 } from 'lucide-react';
import { motion } from 'framer-motion';

const defaultConfig = `version: "1.0"
workspace_dir: "./workspace"

repositories:
  - name: backend-api
    url: https://github.com/cyclic-software/starter-express-api.git
    path: ./services/backend
    setup_commands:
      - npm install

services:
  - name: backend
    repo: backend-api
    run_command: npm start`;

export const ConfigEditor = ({ onUpload, config }) => {
  const [configText, setConfigText] = useState(defaultConfig);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

  const handleUpload = async () => {
    setIsLoading(true);
    setError(null);
    setSuccess(false);

    try {
      await onUpload(configText, (response) => {
        setIsLoading(false);
        if (response.type === 'success' && response.payload.valid) {
          setSuccess(true);
          setTimeout(() => setSuccess(false), 3000);
        } else if (response.type === 'error') {
          setError(response.payload.message);
        } else if (response.payload.errors) {
          setError(response.payload.errors.join('\n'));
        }
      });
    } catch (err) {
      setIsLoading(false);
      setError(err.message);
    }
  };

  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setConfigText(e.target.result);
      };
      reader.readAsText(file);
    }
  };

  return (
    <motion.div
      className="card"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
    >
      <div className="flex items-center justify-between mb-4">
        <div>
          <h2 className="text-xl font-bold text-text-primary flex items-center gap-2">
            <FileText className="w-5 h-5 text-primary-500" />
            Configuration
          </h2>
          <p className="text-sm text-text-tertiary mt-1">
            Upload or edit your YAML configuration
          </p>
        </div>

        {config && (
          <div className="flex items-center gap-2 text-sm">
            <div className="badge-success">
              <CheckCircle className="w-3 h-3" />
              <span>{config.repositories} repos</span>
            </div>
            <div className="badge-info">
              <span>{config.services} services</span>
            </div>
          </div>
        )}
      </div>

      {/* File Upload Button */}
      <div className="mb-4">
        <label className="btn-secondary cursor-pointer inline-flex">
          <Upload className="w-4 h-4" />
          Upload YAML File
          <input
            type="file"
            accept=".yaml,.yml"
            onChange={handleFileUpload}
            className="hidden"
          />
        </label>
      </div>

      {/* Config Editor */}
      <div className="mb-4">
        <textarea
          value={configText}
          onChange={(e) => setConfigText(e.target.value)}
          className="w-full h-64 px-4 py-3 bg-surface-base border border-surface-border rounded-lg text-text-primary font-mono text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent resize-none"
          placeholder="Paste your YAML configuration here..."
          spellCheck={false}
        />
      </div>

      {/* Error Display */}
      {error && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-4 p-4 bg-red-500/10 border border-red-500/20 rounded-lg flex items-start gap-3"
        >
          <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" />
          <div className="flex-1">
            <h4 className="text-red-400 font-medium mb-1">Configuration Error</h4>
            <pre className="text-sm text-red-300 whitespace-pre-wrap font-mono">
              {error}
            </pre>
          </div>
        </motion.div>
      )}

      {/* Success Display */}
      {success && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-4 p-4 bg-green-500/10 border border-green-500/20 rounded-lg flex items-center gap-3"
        >
          <CheckCircle className="w-5 h-5 text-green-400" />
          <span className="text-green-400 font-medium">
            Configuration uploaded successfully!
          </span>
        </motion.div>
      )}

      {/* Actions */}
      <div className="flex items-center gap-3">
        <button
          onClick={handleUpload}
          disabled={isLoading || !configText.trim()}
          className="btn-primary"
        >
          {isLoading ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin" />
              Uploading...
            </>
          ) : (
            <>
              <Upload className="w-4 h-4" />
              Upload Configuration
            </>
          )}
        </button>

        {config && (
          <div className="text-sm text-text-tertiary">
            Workspace: <span className="text-text-secondary font-mono">{config.workspace_dir}</span>
          </div>
        )}
      </div>
    </motion.div>
  );
};
