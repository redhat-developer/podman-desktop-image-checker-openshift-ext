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
  deleteImage,
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
const providerName = 'Red Hat OpenShift Checker';
const extensionName = 'Red Hat OpenShift Checker extension';

test.use({
  runnerOptions: new RunnerOptions({ customFolder: 'image-checker-tests-pd', autoUpdate: false, autoCheckUpdates: false }),
});
test.beforeAll(async ({ runner, page, welcomePage }) => {
  runner.setVideoAndTraceName('image-checker-e2e');
  await welcomePage.handleWelcomePage(true);
  extensionCard = new ExtensionCardPage(page, extensionLabelName, extensionLabel);
});

test.afterAll(async ({ runner, page }) => {
  test.setTimeout(60000);
  deleteImage(page, imageToCheck);
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
    await disableAndRemoveExtension(navigationBar);
  });

  test('Extension can be installed from an OCI image', async ({ navigationBar }) => {
    test.setTimeout(180000);
    const extensions = await navigationBar.openExtensions();
    await extensions.installExtensionFromOCIImage(imageName);
    await playExpect(extensionCard.card).toBeVisible();
  });

  test('Extension card is present in extension list, extension is active', async ({ navigationBar }) => {
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
    await extensionCard.openExtensionDetails(extensionName);
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

test.describe.serial('Red Hat Image Checker extension functionality', () => {
  test('Pull testing image, check tab is present in image details page', async ({ navigationBar }) => {
    const imagesPage = await navigationBar.openImages();
    await playExpect(imagesPage.heading).toBeVisible();

    const pullImagePage = await imagesPage.openPullImage();
    await playExpect(pullImagePage.heading).toBeVisible();
    await pullImagePage.pullImage(imageToCheck);
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).toBeVisible();

    await imageDetailsPage.imageCheckerTab.click();
    await playExpect(imageDetailsPage.imageCheckerTabContent).toBeVisible();
  });

  test('Image checker can be turned off and on in image details page', async ({ navigationBar }) => {
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).toBeVisible();
    await imageDetailsPage.imageCheckerTab.click();

    const provider = await imageDetailsPage.getProvider(providerName);
    await playExpect(provider).toBeVisible();
    const analysisResults = imageDetailsPage.analysisTable;
    await playExpect(analysisResults).toBeVisible();

    await imageDetailsPage.setProviderCheckbox(providerName, false);
    await playExpect(await imageDetailsPage.getProviderCheckbox(providerName)).not.toBeChecked();
    await playExpect(analysisResults.getByRole("row")).not.toBeVisible();

    await imageDetailsPage.setProviderCheckbox(providerName, true);
    await playExpect(await imageDetailsPage.getProviderCheckbox(providerName)).toBeChecked();
    await playExpect(analysisResults.getByRole("row")).toBeVisible();
  });

  test('Verify analysis status, analysis providers and results are visible', async ({ navigationBar }) => {
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).toBeVisible();
    await imageDetailsPage.imageCheckerTab.click();
    await playExpect(imageDetailsPage.imageCheckerTabContent).toBeVisible();

    // wait for the analysis to be complete
    const analysisStatus = imageDetailsPage.analysisStatus;
    await playExpect(analysisStatus).toBeVisible();
    await playExpect
      .poll(async () => await analysisStatus.innerText(), { timeout: 60000 })
      .toContain('Image analysis complete');

    await playExpect(imageDetailsPage.providersTable).toBeVisible();

    const redhatCheckerProvider = await imageDetailsPage.getProvider(providerName);
    await playExpect(redhatCheckerProvider).toBeVisible();

    await playExpect(imageDetailsPage.analysisTable).toBeVisible();
  });

  test('Verify analysis results', async ({ navigationBar }) => {
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).toBeVisible();
    await imageDetailsPage.imageCheckerTab.click();

    if (isLinux) {
      const resultError = await imageDetailsPage.getAnalysisResult('Analyze error');
      await playExpect(resultError, 'Tests assume analysis not implemented on linux').toBeVisible();
    }
    else {
      const checkExpose = await imageDetailsPage.getAnalysisResult('Privileged port exposed');
      await playExpect(checkExpose).toBeVisible();
      const checkChown = await imageDetailsPage.getAnalysisResult('Owner set');
      await playExpect(checkChown).toBeVisible();
      const checkUser = await imageDetailsPage.getAnalysisResult('User set to root');
      await playExpect(checkUser).toBeVisible();
      // TODO: create custom image to test specific directives, instead of using httpd
    }
  });
});

test.describe.serial('Red Hat Image Checker extension handling', () => {
  test('Extension can be disabled', async ({ navigationBar }) => {
    const extensions = await navigationBar.openExtensions();
    await playExpect
    .poll(async () => await extensions.extensionIsInstalled(extensionLabel), { timeout: 15000 })
    .toBeTruthy();    
    const extensionCard = await extensions.getInstalledExtension(extensionLabelName, extensionLabel);
    await extensionCard.disableExtension();
    await playExpect(extensionCard.status).toHaveText(disabledExtensionStatus);
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).not.toBeVisible();
  });

  test('Extension can be re-enabled', async ({ navigationBar }) => {
    const extensions = await navigationBar.openExtensions();
    await playExpect
    .poll(async () => await extensions.extensionIsInstalled(extensionLabel), { timeout: 15000 })
    .toBeTruthy();
    const extensionCard = await extensions.getInstalledExtension(extensionLabelName, extensionLabel);
    await extensionCard.enableExtension();
    await playExpect(extensionCard.status).toHaveText(activeExtensionStatus);
    const imageDetailsPage = await getImageDetailsPage(navigationBar);
    await playExpect(imageDetailsPage.imageCheckerTab).toBeVisible();
  });

  test('Extension can be removed', async ({ navigationBar }) => {
    await disableAndRemoveExtension(navigationBar);
  });
});

async function disableAndRemoveExtension(navigationBar: NavigationBar): Promise<void> {
  const extensions = await navigationBar.openExtensions();
  const extensionCard = await extensions.getInstalledExtension(extensionLabelName, extensionLabel);
  playExpect(extensionCard.status).toHaveText(activeExtensionStatus);
  await playExpect(extensionCard.removeButton).toBeVisible();
  await extensionCard.removeExtension();
  await playExpect
    .poll(async () => await extensions.extensionIsInstalled(extensionLabel), { timeout: 15000 })
    .toBeFalsy();
}

async function getImageDetailsPage(navigationBar: NavigationBar): Promise<ImageCheckerImageDetailsPage> {
  const imagesPage = await navigationBar.openImages();
  await playExpect(imagesPage.heading).toBeVisible();

  const exists = await imagesPage.waitForImageExists(imageToCheck);
  playExpect(exists, `${imageToCheck} image not found in the image list`).toBeTruthy();

  const imageDetailPage = await imagesPage.openImageDetails(imageToCheck);
  const imageDetailsPage = new ImageCheckerImageDetailsPage(imageDetailPage.page, imageToCheck);
  return imageDetailsPage;
}
