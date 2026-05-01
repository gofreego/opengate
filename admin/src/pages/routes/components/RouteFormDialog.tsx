import { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  FormControlLabel,
  Switch,
  Typography,
  IconButton,
  Chip,
  Paper,
  Divider,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  OutlinedInput,
} from '@mui/material'
import { Add as AddIcon, Close as CloseIcon, Delete as DeleteIcon } from '@mui/icons-material'
import type { Config, CreateConfigRequest, UpdateConfigRequest, Authentication, AuthenticationException } from '../../../apis/proto/opengate/v1/config'

const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD']

interface ConfigFormDialogProps {
  open: boolean
  onClose: () => void
  onSave: (data: CreateConfigRequest | UpdateConfigRequest, isEdit: boolean) => Promise<void>
  editData: Config | null
}

export const RouteFormDialog = ({
  open,
  onClose,
  onSave,
  editData,
}: ConfigFormDialogProps) => {
  const [name, setName] = useState('')
  const [pathPrefix, setPathPrefix] = useState('')
  const [targetUrl, setTargetUrl] = useState('')
  const [stripPrefix, setStripPrefix] = useState(false)
  const [authRequired, setAuthRequired] = useState(false)
  const [authExcept, setAuthExcept] = useState<AuthenticationException[]>([])
  const [middleware, setMiddleware] = useState<string[]>([])
  const [newMiddleware, setNewMiddleware] = useState('')
  const [timeout, setTimeout] = useState('30000000000') // 30s in nanoseconds
  const [saving, setSaving] = useState(false)

  // New exception form state
  const [newExceptPath, setNewExceptPath] = useState('')
  const [newExceptMethods, setNewExceptMethods] = useState<string[]>([])
  const [showAddException, setShowAddException] = useState(false)

  useEffect(() => {
    if (editData) {
      setName(editData.name)
      setPathPrefix(editData.pathPrefix)
      setTargetUrl(editData.targetUrl)
      setStripPrefix(editData.stripPrefix)
      setAuthRequired(editData.authentication?.required || false)
      setAuthExcept(editData.authentication?.except || [])
      setMiddleware(editData.middleware || [])
      setTimeout(editData.timeout || '30000000000')
    } else {
      resetForm()
    }
  }, [editData, open])

  const resetForm = () => {
    setName('')
    setPathPrefix('')
    setTargetUrl('')
    setStripPrefix(false)
    setAuthRequired(false)
    setAuthExcept([])
    setMiddleware([])
    setNewMiddleware('')
    setNewExceptPath('')
    setNewExceptMethods([])
    setShowAddException(false)
    setTimeout('30000000000')
  }

  const handleAddMiddleware = () => {
    if (newMiddleware.trim() && !middleware.includes(newMiddleware.trim())) {
      setMiddleware([...middleware, newMiddleware.trim()])
      setNewMiddleware('')
    }
  }

  const handleRemoveMiddleware = (item: string) => {
    setMiddleware(middleware.filter((m) => m !== item))
  }

  const handleAddException = () => {
    if (newExceptPath.trim()) {
      setAuthExcept([...authExcept, { path: newExceptPath.trim(), methods: newExceptMethods }])
      setNewExceptPath('')
      setNewExceptMethods([])
      setShowAddException(false)
    }
  }

  const handleRemoveException = (index: number) => {
    setAuthExcept(authExcept.filter((_, i) => i !== index))
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const authentication: Authentication = {
        required: authRequired,
        except: authExcept,
      }

      const data: CreateConfigRequest | UpdateConfigRequest = {
        name,
        pathPrefix,
        targetUrl,
        stripPrefix,
        authentication,
        middleware,
        timeout,
        ...(editData && { id: editData.id }),
      }

      await onSave(data, !!editData)
      onClose()
    } finally {
      setSaving(false)
    }
  }

  const isValid = name.trim() && pathPrefix.trim() && targetUrl.trim()

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {editData ? 'Edit Config' : 'Create Config'}
      </DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
          <TextField
            label="Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            fullWidth
            required
            placeholder="e.g., user-service"
          />
          <TextField
            label="Path Prefix"
            value={pathPrefix}
            onChange={(e) => setPathPrefix(e.target.value)}
            fullWidth
            required
            placeholder="e.g., /api/users"
          />
          <TextField
            label="Target URL"
            value={targetUrl}
            onChange={(e) => setTargetUrl(e.target.value)}
            fullWidth
            required
            placeholder="e.g., http://user-service:8080"
          />
          <FormControlLabel
            control={
              <Switch
                checked={stripPrefix}
                onChange={(e) => setStripPrefix(e.target.checked)}
              />
            }
            label="Strip Path Prefix"
          />
          
          <Divider sx={{ my: 1 }} />
          
          {/* Authentication Section */}
          <Typography variant="subtitle1" fontWeight={600}>
            Authentication
          </Typography>
          
          <FormControlLabel
            control={
              <Switch
                checked={authRequired}
                onChange={(e) => setAuthRequired(e.target.checked)}
              />
            }
            label="Require Authentication"
          />
          
          {/* Authentication Exceptions */}
          <Box sx={{ ml: 2 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
              <Typography variant="subtitle2">
                Exceptions {authRequired ? '(paths that do NOT require auth)' : '(paths that DO require auth)'}
              </Typography>
              <Button
                variant="outlined"
                size="small"
                startIcon={<AddIcon />}
                onClick={() => setShowAddException(true)}
                disabled={showAddException}
              >
                Add
              </Button>
            </Box>
            
            {/* Add new exception form - only shown when plus is clicked */}
            {showAddException && (
              <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
                <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', alignItems: 'flex-start' }}>
                  <TextField
                    size="small"
                    label="Path"
                    value={newExceptPath}
                    onChange={(e) => setNewExceptPath(e.target.value)}
                    placeholder="/health"
                    sx={{ flex: 1, minWidth: 150 }}
                    autoFocus
                  />
                  <FormControl size="small" sx={{ flex: 1, minWidth: 180 }}>
                    <InputLabel>Methods</InputLabel>
                    <Select<string[]>
                      multiple
                      value={newExceptMethods}
                      onChange={(e) => {
                        const value = e.target.value
                        setNewExceptMethods(typeof value === 'string' ? value.split(',') : value)
                      }}
                      input={<OutlinedInput label="Methods" />}
                      renderValue={(selected) => (
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {selected.length === 0 ? (
                            <Typography variant="body2" color="text.secondary">All methods</Typography>
                          ) : (
                            selected.map((value) => (
                              <Chip key={value} label={value} size="small" />
                            ))
                          )}
                        </Box>
                      )}
                    >
                      {HTTP_METHODS.map((method) => (
                        <MenuItem key={method} value={method}>
                          {method}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                  <Button
                    variant="contained"
                    size="small"
                    onClick={handleAddException}
                    disabled={!newExceptPath.trim()}
                    sx={{ mt: 0.5 }}
                  >
                    Add
                  </Button>
                  <Button
                    variant="outlined"
                    size="small"
                    onClick={() => {
                      setShowAddException(false)
                      setNewExceptPath('')
                      setNewExceptMethods([])
                    }}
                    sx={{ mt: 0.5 }}
                  >
                    Cancel
                  </Button>
                </Box>
              </Paper>
            )}
            
            {/* List of exceptions */}
            {authExcept.length > 0 && (
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                {authExcept.map((exc, index) => (
                  <Paper
                    key={index}
                    variant="outlined"
                    sx={{ 
                      p: 1.5, 
                      display: 'flex', 
                      alignItems: 'center', 
                      justifyContent: 'space-between',
                      bgcolor: 'action.hover',
                    }}
                  >
                    <Box>
                      <Typography variant="body2" fontWeight={500} sx={{ fontFamily: 'monospace' }}>
                        {exc.path}
                      </Typography>
                      <Box sx={{ display: 'flex', gap: 0.5, mt: 0.5, flexWrap: 'wrap' }}>
                        {exc.methods.length > 0 ? (
                          exc.methods.map((method) => (
                            <Chip
                              key={method}
                              label={method}
                              size="small"
                              variant="outlined"
                              color="primary"
                            />
                          ))
                        ) : (
                          <Chip label="ALL METHODS" size="small" variant="outlined" />
                        )}
                      </Box>
                    </Box>
                    <IconButton
                      size="small"
                      color="error"
                      onClick={() => handleRemoveException(index)}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </Paper>
                ))}
              </Box>
            )}
          </Box>
          
          <Divider sx={{ my: 1 }} />
          
          <TextField
            label="Timeout (nanoseconds)"
            value={timeout}
            onChange={(e) => setTimeout(e.target.value)}
            fullWidth
            type="text"
            helperText="Default: 30000000000 (30 seconds)"
          />
          
          {/* Middleware */}
          <Box>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>
              Middleware
            </Typography>
            <Box sx={{ display: 'flex', gap: 1, mb: 1 }}>
              <TextField
                size="small"
                value={newMiddleware}
                onChange={(e) => setNewMiddleware(e.target.value)}
                placeholder="Add middleware..."
                onKeyPress={(e) => e.key === 'Enter' && handleAddMiddleware()}
              />
              <IconButton onClick={handleAddMiddleware} size="small">
                <AddIcon />
              </IconButton>
            </Box>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
              {middleware.map((m) => (
                <Chip
                  key={m}
                  label={m}
                  size="small"
                  onDelete={() => handleRemoveMiddleware(m)}
                  deleteIcon={<CloseIcon />}
                />
              ))}
            </Box>
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          variant="contained"
          onClick={handleSave}
          disabled={!isValid || saving}
        >
          {saving ? 'Saving...' : editData ? 'Update' : 'Create'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
