import * as React from 'react'
import { Box, Typography, useMediaQuery, useTheme } from '@mui/material'

interface PageHeaderProps {
  title: string
  action?: React.ReactNode
  inlineAction?: React.ReactNode
}

const PageHeader: React.FC<PageHeaderProps> = ({ title, action, inlineAction }) => {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'))
  return (
    <Box
      display="flex"
      alignItems="center"
      justifyContent={isMobile ? 'space-between' : 'center'}
      position="relative"
      mb={3}
      pb={1.5}
      sx={{
        borderBottom: '2px solid hsl(var(--primary))',
        flexDirection: isMobile ? 'column' : 'row',
        gap: isMobile ? 1 : 0,
        textAlign: isMobile ? 'center' : undefined,
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap', justifyContent: 'center' }}>
        <Typography
          variant="h5"
          sx={{
            color: 'hsl(var(--text-primary))',
            fontWeight: 700,
            fontFamily: "'Inter', sans-serif",
            letterSpacing: '-0.01em',
            textAlign: 'center',
          }}
        >
          {title}
        </Typography>
        {inlineAction}
      </Box>
      {action && (
        <Box
          sx={{
            position: isMobile ? 'relative' : 'absolute',
            right: isMobile ? 'auto' : 0,
            top: isMobile ? 'auto' : '50%',
            transform: isMobile ? 'none' : 'translateY(-50%)',
            mt: isMobile ? 0.5 : 0,
          }}
        >
          {action}
        </Box>
      )}
    </Box>
  )
}

export default PageHeader
