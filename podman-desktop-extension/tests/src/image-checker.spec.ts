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

import type { NavigationBar } from '@podman-desktop/tests-playwright';
import {
  expect as playExpect,
  ExtensionCardPage,
  RunnerOptions,
  test,
} from '@podman-desktop/tests-playwright';
import { ImageCheckerExtensionPage } from './pages/image-checker-extension-page';
import { ImageCheckerImageDetailsPage } from './pages/image-checker-image-details-page';

let extensionInstalled = false;
let extensionCard: ExtensionCardPage;
const imageName = 'ghcr.io/redhat-developer/podman-desktop-image-checker-openshift-ext:latest';
const extensionLabel = 'redhat.openshift-checker';
const extensionLabelName = 'openshift-checker';
const activeExtensionStatus = 'ACTIVE';
const disabledExtensionStatus = 'DISABLED';
const imageToCheck = 'ghcr.io/linuxcontainers/alpine';
const isLinux = process.platform === 'linux';
test.use({
  runnerOptions: new RunnerOptions({ customFolder: 'image-checker-tests-pd', autoUpdate: false, autoCheckUpdates: false }),
});
test.beforeAll(async ({ runner, page, welcomePage }) => {
  runner.setVideoAndTraceName('image-checker-e2e');
  await welcomePage.handleWelcomePage(true);
  extensionCard = new ExtensionCardPage(page, extensionLabelName, extensionLabel);
});

test.afterAll(async ({ runner }) => {
  await runner.close();
});

test.describe.serial('Red Hat Image Checker extension installation', () => {
  test('Check if extension is already installed', async ({ navigationBar }) => {
    const extensions = await navigationBar.openExtensions();
    if (await extensions.extensionIsInstalled(extensionLabel)) {
      extensionInstalled = true;
    }
  });

  test('Remove old version of the extension', async ({ navigationBar }) => {
    test.skip(!extensionInstalled);
    await removeExtension(navigationBar);
  });

  test('Extension can be installed from an OCI image', async ({ navigationBar }) => {
    test.setTimeout(180000);
    const extensions = await navigationBar.openExtensions();
    await extensions.installExtensionFromOCIImage(imageName);
    await playExpect(extensionCard.card).toBeVisible();
  });

  test('Extension is installed and active, extension card is present', async ({ navigationBar }) => {
    const extensions = await navigationBar.openExtensions();
    await playExpect
      .poll(async () => await extensions.extensionIsInstalled(extensionLabel), { timeout: 60000 })
      .toBeTruthy();
    const extensionCard = await extensions.getInstalledExtension(extensionLabelName, extensionLabel);
    await playExpect(extensionCard.status, `Extension status is: ${extensionCard.status}`).toHaveText(activeExtensionStatus);
  });

  test("Extension details show correct status, no error", async ({ page, navigationBar }) => {
    const extensions = await navigationBar.openExtensions();
    const extensionCard = await extensions.getInstalledExtension(extensionLabelName, extensionLabel);
    await extensionCard.openExtensionDetails('Red Hat OpenShift Checker extension');
    const detailsPage = new ImageCheckerExtensionPage(page);
    await playExpect(detailsPage.heading).toBeVisible();
    await playExpect(detailsPage.status).toHaveText(activeExtensionStatus);
    const errorTab = detailsPage.tabs.getByRole('button', { name: 'Error' });
    // we would like to propagate the error's stack trace into test failure message
    let stackTrace = '';
    if ((await errorTab.count()) > 0) {
      await detailsPage.activateTab('Error');
      stackTrace = await detailsPage.errorStackTrace.innerText();
    }
    await playExpect(errorTab, `Error Tab was present with stackTrace: ${stackTrace}`).not.toBeVisible();
  });
});

test.describe.serial('Red Hat Image Checker extension verification', () => {
  test('Extension is installed, check tab present in image details page', async ({ navigationBar }) => {
    const extensions = await navigationBar.openExtensions();
    if (await extensions.extensionIsInstalled(extensionLabel)) {
      extensionInstalled = true;
    }
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).toBeVisible();
    await imageDetailsPage.imageCheckerTab.click();
    await playExpect(imageDetailsPage.imageCheckerTabContent).toBeVisible();
  });

  test('Image checker can be turned off and on in image details page', async ({ navigationBar }) => {
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).toBeVisible();
    await imageDetailsPage.imageCheckerTab.click();

    const provider = await imageDetailsPage.getProvider('Red Hat OpenShift Checker');
    await playExpect(provider).toBeVisible();
    const analysisResults = imageDetailsPage.analysisTable;
    await playExpect(analysisResults).toBeVisible();
    
    // check if analysis is hidden
    const providerCheckbox = await imageDetailsPage.getProviderCheckbox(provider);
    await providerCheckbox.click();
    await playExpect(analysisResults.getByRole("row")).not.toBeVisible();

    // check if analysis is visible
    await providerCheckbox.click();
    await playExpect(analysisResults.getByRole("row")).toBeVisible();
  });

  test('Checker is present in image details page, analysis is visible', async ({ navigationBar }) => {
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await imageDetailsPage.imageCheckerTab.click();
    await playExpect(imageDetailsPage.imageCheckerTabContent).toBeVisible();

    // wait for the analysis to be complete
    const analysisStatus = await imageDetailsPage.getAnalysisStatus();
    await playExpect(analysisStatus).toBeVisible();
    await playExpect
    .poll(async () => await analysisStatus.innerText(), { timeout: 60000 })
    .toContain('Image analysis complete');

    // test the providers table
    const providersTable = imageDetailsPage.providersTable;
    await playExpect(providersTable).toBeVisible();
    const redhatCheckerProvider = await imageDetailsPage.getProvider('Red Hat OpenShift Checker');
    await playExpect(redhatCheckerProvider).toBeVisible();

    const analysisTable = imageDetailsPage.analysisTable;
    await playExpect(analysisTable).toBeVisible();

    // test all the directives
    if (isLinux) {
      const resultError = await imageDetailsPage.getAnalysisResult('Analyze error');
      await playExpect(resultError).toBeVisible();
    }
    else {
      const resultRunDir = await imageDetailsPage.getAnalysisResult('Priviliged port exposed');
      await playExpect(resultRunDir).toBeVisible();
      const resultOwnDir = await imageDetailsPage.getAnalysisResult('Owner set');
      await playExpect(resultOwnDir).toBeVisible();
      const resultUserDir = await imageDetailsPage.getAnalysisResult('User set to root');
      await playExpect(resultUserDir).toBeVisible();
      // TODO: create custom image to test specific directives, instead of using httpd
    }
  });
});

test.describe.serial('Red Hat Image Checker extension handling', () => {
  test('Extension can be disabled', async ({ navigationBar }) => {
    const extensions = await navigationBar.openExtensions();
    playExpect(extensions.extensionIsInstalled(extensionLabel)).toBeTruthy();
    const extensionCard = await extensions.getInstalledExtension(extensionLabelName, extensionLabel);
    await extensionCard.disableExtension();
    await playExpect(extensionCard.status).toHaveText(disabledExtensionStatus);
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).not.toBeVisible();
  });

  test('Extension can be re-enabled', async ({ navigationBar }) => {
    const extensions = await navigationBar.openExtensions();
    playExpect(extensions.extensionIsInstalled(extensionLabel)).toBeTruthy();
    const extensionCard = await extensions.getInstalledExtension(extensionLabelName, extensionLabel);
    await extensionCard.enableExtension();
    await playExpect(extensionCard.status).toHaveText(activeExtensionStatus);
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).toBeVisible();
  });

  test('Extension can be removed', async ({ navigationBar }) => {
    await removeExtension(navigationBar);
  });
});

async function removeExtension(navigationBar: NavigationBar): Promise<void> {
  const extensions = await navigationBar.openExtensions();
  const extensionCard = await extensions.getInstalledExtension(extensionLabelName, extensionLabel);
  await extensionCard.disableExtension();
  await extensionCard.removeExtension();
  await playExpect
    .poll(async () => await extensions.extensionIsInstalled(extensionLabel), { timeout: 15000 })
    .toBeFalsy();
}

async function getImageDetailsPage(navigationBar: NavigationBar): Promise<ImageCheckerImageDetailsPage> {
  const imagesPage = await navigationBar.openImages();
  await playExpect(imagesPage.heading).toBeVisible();

  const pullImagePage = await imagesPage.openPullImage();
  const updatedImages = await pullImagePage.pullImage(imageToCheck);

  const exists = await updatedImages.waitForImageExists(imageToCheck);
  playExpect(exists, `${imageToCheck} image not found in the image list`).toBeTruthy();

  const imageDetailPage = await imagesPage.openImageDetails(imageToCheck);
  const imageDetailsPage = new ImageCheckerImageDetailsPage(imageDetailPage.page, imageToCheck);
  return imageDetailsPage;
}
