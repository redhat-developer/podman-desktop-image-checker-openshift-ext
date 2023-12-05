/**********************************************************************
 * Copyright (C) 2022 Red Hat, Inc.
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
import * as extensionApi from '@podman-desktop/api';

// Activate the extension asynchronously
export async function activate(extensionContext: extensionApi.ExtensionContext): Promise<void> {
  // Create a provider with an example name, ID and icon
  const provider = extensionApi.provider.createProvider({
    name: 'OpenShift Checker',
    id: 'openshift-checker',
    status: 'unknown',
    images: {
      icon: './icon.png',
      logo: './icon.png',
    },
  });
  extensionContext.subscriptions.push(provider);

  const checker = extensionApi.imageChecker.registerImageCheckerProvider({
    check: async (
      image: extensionApi.ImageInfo,
      token?: extensionApi.CancellationToken,
    ): Promise<extensionApi.ImageChecks> => {
      let path: string;
      if (extensionApi.env.isLinux) {
        path = `${__dirname}/../doa.linux`;
      } else if (extensionApi.env.isMac) {
        path = `${__dirname}/../doa.darwin`;
      } else if (extensionApi.env.isWindows) {
        path = `${__dirname}\\..\\doa.exe`;
      }
      const result = await extensionApi.process.exec(
        path,
        ['analyze', '-i', image.Id, '-o', 'json'],
        { 
          token,
        }
      );

      const originalChecks = JSON.parse(result.stdout);

      return {
        checks: originalChecks.map(c => {
          return {
            name: c.name,
            status: c.status,
            markdownDescription: c.description,
            severity: c.severity,
          } as extensionApi.ImageCheck;
        }),
      };
    },
  });
  extensionContext.subscriptions.push(checker);
}

// Deactivate the extension
export function deactivate(): void {
  console.log('stopping openshift-checker extension');
}
