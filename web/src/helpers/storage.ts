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

interface StorageShape {
  getItem: (key: string) => string | null;
  removeItem: (key: string) => void;
  setItem: (key: string, item: string) => void;
  clear: () => void;
}

class Storage {
  storageInstance: StorageShape;

  constructor(storageInstance: StorageShape) {
    this.storageInstance = storageInstance;
  }

  /**
   * Stores data in the storage
   *
   * @param key The key to store the data under
   * @param data The data to store
   */
  write(key: string, data: any) {
    const storedShape = data instanceof Object ? JSON.stringify(data) : data;
    this.storageInstance.setItem(key, storedShape);
  }

  /**
   * Retrieves data from the storage
   *
   * @param key The key to read
   * @returns the data matching this key
   */
  read<T = string>(key: string): T {
    const data = this.storageInstance.getItem(key);
    if (data === null) {
      return null;
    }

    try {
      return JSON.parse(this.storageInstance.getItem(key)) as T;
    } catch (e) {
      return (data as unknown) as T;
    }
  }

  /**
   *
   * @param key The key to delete
   * @returns void
   */
  delete(key: string) {
    this.storageInstance.removeItem(key);
  }

  /**
   * Clears the storage from all of its keys
   */
  clear() {
    this.storageInstance.clear();
  }
}

const storage = new Storage(localStorage);

export default storage;
