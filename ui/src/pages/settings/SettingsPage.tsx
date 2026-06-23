import {
  Container,
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Switch,
  FormControlLabel,
  CircularProgress,
  Alert,
  Divider,
} from '@mui/material'
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'
import SecurityIcon from '@mui/icons-material/Security'
import SettingsIcon from '@mui/icons-material/Settings'
import { useState } from 'react'
import { PageHeader } from '../../components'
import { settingsService, type CORSConfig } from '../../services'

const DEFAULT_CORS: CORSConfig = {
  enabled: true,
  allowedOrigins: '*',
  allowedMethods: 'GET, POST, PUT, DELETE, OPTIONS, PATCH',
  allowedHeaders: 'Accept, Authorization, Content-Type, X-CSRF-Token, X-User-Id, X-User-Perms',
  maxAge: 3600,
}

function CORSSection() {
  const [expanded, setExpanded] = useState(false)
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)
  const [config, setConfig] = useState<CORSConfig | null>(null)

  const handleExpand = async (_: React.SyntheticEvent, isExpanded: boolean) => {
    setExpanded(isExpanded)
    if (isExpanded && config === null) {
      setLoading(true)
      setError(null)
      try {
        const fetched = await settingsService.getCORSConfig()
        setConfig(fetched ?? DEFAULT_CORS)
      } catch (e: unknown) {
        setError(e instanceof Error ? e.message : 'Failed to load CORS config')
        setConfig(DEFAULT_CORS)
      } finally {
        setLoading(false)
      }
    }
  }

  const handleSave = async () => {
    if (!config) return
    setSaving(true)
    setError(null)
    setSuccess(false)
    try {
      await settingsService.saveCORSConfig(config)
      setSuccess(true)
      setTimeout(() => setSuccess(false), 3000)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to save CORS config')
    } finally {
      setSaving(false)
    }
  }

  const update = (field: keyof CORSConfig, value: unknown) =>
    setConfig((prev) => (prev ? { ...prev, [field]: value } : prev))

  return (
    <Accordion expanded={expanded} onChange={handleExpand}>
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <SecurityIcon fontSize="small" color="primary" />
          <Typography variant="subtitle1" fontWeight={500}>
            CORS Configuration
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ ml: 1 }}>
            — Cross-Origin Resource Sharing policy for all servers
          </Typography>
        </Box>
      </AccordionSummary>

      <AccordionDetails>
        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 3 }}>
            <CircularProgress size={28} />
          </Box>
        ) : (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            {error && <Alert severity="error">{error}</Alert>}
            {success && <Alert severity="success">CORS configuration saved.</Alert>}

            {config && (
              <>
                <FormControlLabel
                  control={
                    <Switch
                      checked={config.enabled}
                      onChange={(e) => update('enabled', e.target.checked)}
                    />
                  }
                  label="Enable CORS"
                />

                <TextField
                  label="Allowed Origins"
                  value={config.allowedOrigins}
                  onChange={(e) => update('allowedOrigins', e.target.value)}
                  helperText='Use * to allow all origins, or specify a single origin (e.g. https://example.com)'
                  fullWidth
                />

                <TextField
                  label="Allowed Methods"
                  value={config.allowedMethods}
                  onChange={(e) => update('allowedMethods', e.target.value)}
                  helperText="Comma-separated HTTP methods (e.g. GET, POST, PUT, DELETE, OPTIONS, PATCH)"
                  fullWidth
                />

                <TextField
                  label="Allowed Headers"
                  value={config.allowedHeaders}
                  onChange={(e) => update('allowedHeaders', e.target.value)}
                  helperText="Comma-separated request headers the client is allowed to send"
                  fullWidth
                />

                <TextField
                  label="Max Age (seconds)"
                  type="number"
                  value={config.maxAge}
                  onChange={(e) => update('maxAge', parseInt(e.target.value, 10) || 0)}
                  helperText="How long the preflight response can be cached by the browser"
                  sx={{ maxWidth: 240 }}
                />

                <Box>
                  <Button
                    variant="contained"
                    onClick={handleSave}
                    disabled={saving}
                    startIcon={saving ? <CircularProgress size={16} /> : undefined}
                  >
                    {saving ? 'Saving…' : 'Save CORS Config'}
                  </Button>
                </Box>
              </>
            )}
          </Box>
        )}
      </AccordionDetails>
    </Accordion>
  )
}

export const SettingsPage = () => {
  const [apiBaseUrl, setApiBaseUrl] = useState(
    localStorage.getItem('opengate_api_base_url') || ''
  )

  const handleSaveApi = () => {
    localStorage.setItem('opengate_api_base_url', apiBaseUrl)
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <PageHeader
        title="Settings"
        subtitle="Configure your OpenGate admin settings"
      />

      {/* API Configuration */}
      <Card sx={{ borderRadius: 2, mb: 3 }}>
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
            <SettingsIcon fontSize="small" color="primary" />
            <Typography variant="h6">API Configuration</Typography>
          </Box>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="API Base URL"
              value={apiBaseUrl}
              onChange={(e) => setApiBaseUrl(e.target.value)}
              placeholder="http://localhost:8080"
              helperText="Leave empty to use the default (same origin)"
              fullWidth
            />
            <Box>
              <Button variant="contained" onClick={handleSaveApi}>
                Save Settings
              </Button>
            </Box>
          </Box>
        </CardContent>
      </Card>

      {/* Gateway Settings — expandable sections */}
      <Typography variant="h6" sx={{ mb: 1 }}>
        Gateway Settings
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Expand a section to view and edit its configuration. Changes are stored in the
        repository and applied at the next backend refresh interval.
      </Typography>

      <Card sx={{ borderRadius: 2 }}>
        <CardContent sx={{ p: 0, '&:last-child': { pb: 0 } }}>
          <CORSSection />
          <Divider />
          {/* Future gateway setting sections go here */}
        </CardContent>
      </Card>
    </Container>
  )
}
