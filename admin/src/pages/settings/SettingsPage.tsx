import { Container, Box, Card, CardContent, Typography, TextField, Button } from '@mui/material'
import { useState } from 'react'
import { PageHeader } from '../../components'

export const SettingsPage = () => {
  const [apiBaseUrl, setApiBaseUrl] = useState(
    localStorage.getItem('opengate_api_base_url') || ''
  )

  const handleSave = () => {
    localStorage.setItem('opengate_api_base_url', apiBaseUrl)
    // Could reload the page or update the httpClient
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <PageHeader 
        title="Settings" 
        subtitle="Configure your OpenGate admin settings"
      />

      <Card sx={{ borderRadius: 2 }}>
        <CardContent>
          <Typography variant="h6" sx={{ mb: 2 }}>
            API Configuration
          </Typography>
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
              <Button variant="contained" onClick={handleSave}>
                Save Settings
              </Button>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Container>
  )
}
