/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
