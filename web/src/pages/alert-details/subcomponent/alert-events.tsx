import React, { useState } from 'react';
import JsonViewer from 'Components/json-viewer';
import Panel from 'Components/panel';
import PaginationControls from 'Components/utils/table-pagination-controls';

interface AlertEventsProps {
  events: string[];
}

const AlertEvents: React.FC<AlertEventsProps> = ({ events }) => {
  // because we are going to use that in PaginationControls we are starting an indexing starting
  // from 1 instead of 0. That's why we are using `eventIndex - 1` when selecting the proper event.
  // Normally the `PaginationControls` are used for displaying pages so they are built with a
  // 1-based indexing in mind
  const [eventIndex, setEventIndex] = useState(1);
  return (
    <Panel
      size="large"
      title="Triggered Events"
      actions={
        <PaginationControls
          page={eventIndex}
          totalPages={events.length}
          onPageChange={setEventIndex}
        />
      }
    >
      <JsonViewer data={JSON.parse(JSON.parse(events[eventIndex - 1]))} />
    </Panel>
  );
};

export default AlertEvents;
