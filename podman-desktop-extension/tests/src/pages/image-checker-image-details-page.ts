/**********************************************************************
 * Copyright (C) 2025 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 ***********************************************************************/

import type { Locator, Page } from '@playwright/test';
import { DetailsPage } from '@podman-desktop/tests-playwright';

export class ImageCheckerImageDetailsPage extends DetailsPage {
  readonly imageCheckerTab: Locator;
  readonly imageCheckerTabContent: Locator;
  readonly providersTable: Locator;
  readonly analysisTable: Locator;

  constructor(page: Page, name: string) {
    super(page, name);
    this.imageCheckerTab = this.tabs.getByText('Check');
    this.imageCheckerTabContent = this.page.getByRole('region', { name: 'Tab Content' });
    this.providersTable = this.imageCheckerTabContent.getByLabel('Providers', { exact: true });
    this.analysisTable = this.imageCheckerTabContent.getByLabel('Analysis Results', { exact: true });
  }

  async getAnalysisStatus(): Promise<Locator> {
    return this.imageCheckerTabContent.getByRole('status', { name: 'Analysis Status' });
  }

  async getProvider(providerName: string): Promise<Locator> {
    return this.providersTable.getByRole('row', {name: providerName});
  }

  async getAnalysisResult(analysisName: string): Promise<Locator> {
    return this.analysisTable.getByRole('row', {name: analysisName});
  }

  async getProviderCheckbox(provider: Locator): Promise<Locator> {
    return provider.getByRole('checkbox').locator('..'); // parent element of checkbox
  }
}
