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

import React from 'react';
import { Flex, Grid, Box, Text, Heading, Label } from 'pouncejs';
import PantherLogoWhite from 'Assets/panther-icon--white.svg';
import { Link } from 'react-router-dom';

interface AuthPageContainerComposition {
  Caption: React.FC<{ title: string; subtitle?: string }>;
  AltOptions: React.FC;
}

interface AuthPageContainer {
  banner?: string;
}

const AuthPageContainer: React.FC<AuthPageContainer> & AuthPageContainerComposition = ({
  children,
  banner,
}) => {
  return (
    <Grid gridTemplateColumns="repeat(3, 1fr)" height="100vh">
      <Box gridColumn="1/2" position="relative">
        <img
          src={banner}
          alt="Generic security illustrations"
          width="100%"
          height="100%"
          style={{ objectFit: 'cover' }}
        />
        <Flex
          position="absolute"
          top="0"
          left="0"
          alignItems="center"
          py={150}
          width="100%"
          height="100%"
          flexDirection="column"
        >
          <Link to="/">
            <img src={PantherLogoWhite} alt="Panther Logo" width="54" height="54" />
          </Link>
          <Flex flexDirection="column" alignItems="center" justifyContent="center" m="auto">
            <Label size="medium" mb={5} color="white" textAlign="center">
              Panther Cloud-Native Cybersecurity SIEM
            </Label>
            <Heading
              size="medium"
              color="white"
              lineHeight={5}
              letterSpacing="0.5px"
              textAlign="center"
            >
              Identify threats, misconfigurations, and suspicious activity in all of your
              environments.
            </Heading>
            <Text size="large" color="white" mt={5} textAlign="center">
              Designed for businesses of any size
            </Text>
          </Flex>
        </Flex>
      </Box>
      <Flex gridColumn="2/4" justifyContent="center" alignItems="center">
        <Box width={460}>{children}</Box>
      </Flex>
    </Grid>
  );
};

/**
 * A compound component for the core caption of this auth page
 */
const AuthPageContainerCaption: AuthPageContainerComposition['Caption'] = ({ title, subtitle }) => (
  <Box mb={8}>
    <Heading is="h1" size="medium" color="grey400">
      {title}
    </Heading>
    {subtitle && (
      <Text is="p" size="large" color="grey200" mt={2}>
        {subtitle}
      </Text>
    )}
  </Box>
);

/**
 * A compounet component to act as a wrapper for any alternative options that the page can have
 */
const AuthPageContainerAlt: AuthPageContainerComposition['AltOptions'] = ({ children }) => (
  <Box position="absolute" right={10} top={10}>
    {children}
  </Box>
);

AuthPageContainer.Caption = AuthPageContainerCaption;
AuthPageContainer.AltOptions = AuthPageContainerAlt;

export default AuthPageContainer;