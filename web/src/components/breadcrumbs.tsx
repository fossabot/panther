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

import { Breadcrumbs as PounceBreadcrumbs } from 'pouncejs';
import * as React from 'react';
import { isGuid, capitalize } from 'Helpers/utils';
import { Link } from 'react-router-dom';
import useRouter from 'Hooks/useRouter';

// @HELP_WANTED
// resource names can get super big. We wanna guard against that. Normally we would just say
// that the name should be truncated to the width it fits in, but there is a problem in our case.
// What happens if we go to "/policies/<BIG_TEXT>/edit"? You would be seeing
// "/policies/<BIG_TEXT....>" and you would never be able to see that last part of the breadcrumbs
// that contains the word "edit". To guard against that, we are saying that the biggest breacrumb
// will be a fixed amount of px so that there is *always* space for some other stuff. This is a
// *hardcoded behaviour* meant to guard us against the resource & policy details pages and the number
// assigned to `maxWidth` is special so thaat it can cover our possible breadcrumb combinations when
// a breadcrumb contains a resourceID or a policyID within it. I can't think of any other solution
// that can fit our usecase that doesn't involve complex JS calculations, so please help out
const widthSentinelStyles = {
  display: 'block',
  maxWidth: '700px',
  whiteSpace: 'nowrap' as const,
  overflow: 'hidden' as const,
  textOverflow: 'ellipsis' as const,
};

const Breadcrumbs: React.FC = () => {
  const {
    location: { pathname },
  } = useRouter();

  const fragments = React.useMemo(() => {
    // split by slash and remove empty-splitted values caused by trailing slashes. We also don't
    // want to display the UUIDs as part of the breadcrumbs (which unfortunately exist in the URL)
    const pathKeys = pathname.split('/').filter(fragment => !!fragment && !isGuid(fragment));

    // return the label (what to show) and the uri of each fragment. The URI is constructed by
    // taking the existing path and removing whatever is after each pathKey (only keeping whatever
    // is before-and-including our key). The key is essentially the URL path itself just prettified
    // for displat
    return pathKeys.map(key => ({
      text: capitalize(decodeURIComponent(key).replace(/-_/g, ' ')),
      href: `${pathname.substr(0, pathname.indexOf(`/${key}/`))}/${key}/`,
    }));
  }, [pathname]);

  return (
    <PounceBreadcrumbs
      items={fragments}
      itemRenderer={item => (
        <Link to={item.href} style={widthSentinelStyles}>
          {item.text}
        </Link>
      )}
    />
  );
};

export default Breadcrumbs;
