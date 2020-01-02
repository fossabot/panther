import { Badge, Box, Grid, Label, Text } from 'pouncejs';
import { Link } from 'react-router-dom';
import urls from 'Source/urls';
import React from 'react';
import { AlertDetails } from 'Generated/schema';
import Linkify from 'linkifyjs/react';
import { SEVERITY_COLOR_MAP } from 'Source/constants';
import { formatDatetime } from 'Helpers/utils';
import Panel from 'Components/panel';

interface AlertDetailsInfoProps {
  alert: AlertDetails;
}

const AlertDetailsInfo: React.FC<AlertDetailsInfoProps> = ({ alert }) => {
  return (
    <Panel size="large" title="Alert Details">
      <Grid gridTemplateColumns="repeat(3, 1fr)" gridGap={6}>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            ID
          </Label>
          <Text size="medium" color="black">
            {alert.alertId}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            RULE ORIGIN
          </Label>
          <Text size="medium" color="black">
            {(
              <Link to={urls.rules.details(alert.rule.id)}>
                {alert.rule.displayName || alert.rule.id}
              </Link>
            ) || 'No rule found'}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            LOG TYPES
          </Label>
          {alert.rule.logTypes.length ? (
            alert.rule.logTypes.map(logType => (
              <Text size="medium" color="black" key={logType}>
                {logType}
              </Text>
            ))
          ) : (
            <Text size="medium" color="black">
              All logs
            </Text>
          )}
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            DESCRIPTION
          </Label>
          <Text size="medium" color={alert.rule.description ? 'black' : 'grey200'}>
            <React.Suspense fallback={<span>{alert.rule.description}</span>}>
              <Linkify>{alert.rule.description || 'No description available'}</Linkify>
            </React.Suspense>
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            RUNBOOK
          </Label>
          <Text size="medium" color={alert.rule.runbook ? 'black' : 'grey200'}>
            <React.Suspense fallback={<span>{alert.rule.runbook}</span>}>
              <Linkify>{alert.rule.runbook || 'No runbook available'}</Linkify>
            </React.Suspense>
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            SEVERITY
          </Label>
          <Badge color={SEVERITY_COLOR_MAP[alert.rule.severity]}>{alert.rule.severity}</Badge>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            TAGS
          </Label>
          {alert.rule.tags.length ? (
            alert.rule.tags.map((tag, index) => (
              <Text size="medium" color="black" key={tag} is="span">
                {tag}
                {index !== alert.rule.tags.length - 1 ? ', ' : null}
              </Text>
            ))
          ) : (
            <Text size="medium" color="grey200">
              No tags assigned
            </Text>
          )}
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            CREATED AT
          </Label>
          <Text size="medium" color="black">
            {formatDatetime(alert.creationTime)}
          </Text>
        </Box>
      </Grid>
    </Panel>
  );
};

export default AlertDetailsInfo;
