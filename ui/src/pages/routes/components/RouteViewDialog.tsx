import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  Chip,
  Divider,
  Paper,
} from '@mui/material'
import type { Config } from '../../../apis/proto/opengate/v1/config'

interface ConfigViewDialogProps {
  open: boolean
  onClose: () => void
  config: Config | null
  onEdit: () => void
}

const formatTimestamp = (timestamp: string | undefined): string => {
  if (!timestamp) return '-'
  const date = new Date(parseInt(timestamp, 10) * 1000)
  return date.toLocaleString()
}

const formatTimeout = (timeout: string | undefined): string => {
  if (!timeout) return '30s'
  const ns = parseInt(timeout, 10)
  const seconds = ns / 1_000_000_000
  return `${seconds}s`
}

export const RouteViewDialog = ({
  open,
  onClose,
  config,
  onEdit,
}: ConfigViewDialogProps) => {
  if (!config) return null

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6">{config.name}</Typography>
          <Button variant="outlined" size="small" onClick={onEdit}>
            Edit
          </Button>
        </Box>
      </DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Box>
            <Typography variant="caption" color="text.secondary">
              ID
            </Typography>
            <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
              {config.id}
            </Typography>
          </Box>

          <Divider />

          <Box>
            <Typography variant="caption" color="text.secondary">
              Path Prefix
            </Typography>
            <Typography variant="body1" sx={{ fontFamily: 'monospace' }}>
              {config.pathPrefix}
            </Typography>
          </Box>

          <Box>
            <Typography variant="caption" color="text.secondary">
              Target URL
            </Typography>
            <Typography variant="body1" sx={{ fontFamily: 'monospace' }}>
              {config.targetUrl}
            </Typography>
          </Box>

          <Box sx={{ display: 'flex', gap: 2 }}>
            <Box>
              <Typography variant="caption" color="text.secondary">
                Strip Prefix
              </Typography>
              <Box>
                <Chip
                  label={config.stripPrefix ? 'Yes' : 'No'}
                  size="small"
                  color={config.stripPrefix ? 'primary' : 'default'}
                />
              </Box>
            </Box>

            <Box>
              <Typography variant="caption" color="text.secondary">
                Authentication
              </Typography>
              <Box>
                <Chip
                  label={config.authentication?.required ? 'Required' : 'None'}
                  size="small"
                  color={config.authentication?.required ? 'warning' : 'default'}
                />
              </Box>
            </Box>

            <Box>
              <Typography variant="caption" color="text.secondary">
                Timeout
              </Typography>
              <Typography variant="body2">
                {formatTimeout(config.timeout)}
              </Typography>
            </Box>
          </Box>

          {/* Authentication Exceptions */}
          {config.authentication?.except && config.authentication.except.length > 0 && (
            <Box>
              <Typography variant="caption" color="text.secondary">
                Auth Exceptions {config.authentication.required ? '(no auth required)' : '(auth required)'}
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mt: 0.5 }}>
                {config.authentication.except.map((exc, index) => (
                  <Paper
                    key={index}
                    variant="outlined"
                    sx={{ p: 1, bgcolor: 'action.hover' }}
                  >
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
                  </Paper>
                ))}
              </Box>
            </Box>
          )}

          {config.middleware && config.middleware.length > 0 && (
            <Box>
              <Typography variant="caption" color="text.secondary">
                Middleware
              </Typography>
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 0.5 }}>
                {config.middleware.map((m) => (
                  <Chip key={m} label={m} size="small" variant="outlined" />
                ))}
              </Box>
            </Box>
          )}

          <Divider />

          <Box sx={{ display: 'flex', gap: 4 }}>
            <Box>
              <Typography variant="caption" color="text.secondary">
                Created At
              </Typography>
              <Typography variant="body2">
                {formatTimestamp(config.createdAt)}
              </Typography>
            </Box>
            <Box>
              <Typography variant="caption" color="text.secondary">
                Updated At
              </Typography>
              <Typography variant="body2">
                {formatTimestamp(config.updatedAt)}
              </Typography>
            </Box>
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  )
}
