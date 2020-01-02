import React from 'react';
import { Checkbox, TableProps } from 'pouncejs';

export interface UseSelectableTableRowsProps<T> {
  /**
   * A list of items that are going to be showcased by the Table. TableItem extends the basic JS
   * object, thus the shape of these items can by anything. Usually they keep the same
   * shape as the one that was returned from the API.
   */
  items: TableProps<T>['items'];

  /**
   * A list of column object that describe each column. More info on the shape of these objects
   * follows down below
   * */
  columns: TableProps<T>['columns'];
}

/**
 * A variation of the table where a first column is added in order to show the serial number of
 * each row
 * */
function useSelectableTableRows<ItemShape>({
  columns,
  items,
}: UseSelectableTableRowsProps<ItemShape>) {
  const [selectedItems, setSelectedItems] = React.useState<
    UseSelectableTableRowsProps<ItemShape>['items']
  >([]);

  /* eslint-disable react/display-name */
  const selectableColumns: TableProps<ItemShape>['columns'] = [
    {
      key: 'selection',
      flex: '0 1 auto',
      renderColumnHeader: () => (
        <Checkbox
          checked={selectedItems.length === items.length}
          onChange={checked => setSelectedItems(checked ? items : [])}
        />
      ),
      renderCell: item => (
        <Checkbox
          checked={selectedItems.includes(item)}
          onChange={(checked, e) => {
            e.stopPropagation();

            setSelectedItems(
              checked
                ? [...selectedItems, item]
                : selectedItems.filter(selectedItem => selectedItem !== item)
            );
          }}
        />
      ),
    },
    ...columns,
  ];
  /* eslint-enable react/display-name */

  return React.useMemo(() => ({ selectableColumns, selectedItems }), [
    items,
    columns,
    selectedItems,
  ]);
}

export default useSelectableTableRows;
