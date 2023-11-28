-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Nov 28, 2023 at 02:54 PM
-- Server version: 10.4.27-MariaDB
-- PHP Version: 8.2.0

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `galaxy`
--

-- --------------------------------------------------------

--
-- Table structure for table `accounts`
--

CREATE TABLE `accounts` (
  `id` int(11) NOT NULL,
  `accountNumber` char(36) NOT NULL,
  `bankNumber` char(36) NOT NULL,
  `accountHolderName` text NOT NULL,
  `accountBalance` float NOT NULL,
  `overdraft` float NOT NULL,
  `availableBalance` float NOT NULL,
  `timestamp` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `accounts`
--

INSERT INTO `accounts` (`id`, `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance`, `timestamp`) VALUES
(1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'Doe,John', 4299.33, 0, 4299.33, 1701073233),
(2, 'a6392391-6d0e-4156-9a5c-d756e4c9347e', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'martins,segun', 1499.97, 0, 1499.97, 1701041261),
(3, '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'fizere,dexter', 5699.9, 0, 5699.9, 1701073233),
(4, 'd39f39de-4625-48f8-ab51-6b4c2c0d4c32', '5206927247', 'Ebite,Zion', 100, 0, 100, 1701173609);

-- --------------------------------------------------------

--
-- Table structure for table `accounts_auth`
--

CREATE TABLE `accounts_auth` (
  `id` int(11) NOT NULL,
  `accountNumber` char(36) NOT NULL,
  `password` varchar(255) NOT NULL,
  `timestamp` int(11) NOT NULL,
  `role` varchar(50) NOT NULL DEFAULT 'guest'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `accounts_auth`
--

INSERT INTO `accounts_auth` (`id`, `accountNumber`, `password`, `timestamp`, `role`) VALUES
(1, '01928918', '7853876e75d306ccce5536afe5a2e6412c614e1b2116bd823fb79b9c3fbee998ce389b61ffb90d2e99ca3cf2ec192a4c0b7efb6409b01b37888fc6e804f7be83', 1698173835, ''),
(2, '123456', '7446df0af45f104c8dd1793d7b928436bc335fafbbb1a4a24886b86443ee1da7404172c54cf22511a897389d0c12a784de809f77579b9bce18f89664b804305b', 1698225647, ''),
(3, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a917d01789b58dfd3a702c715496269886f5d363d7445f42ee7b963e9de2a1da7dfbf0b88248ca648e69927353c0a76aaccd1d9b2ef1e32a7fe18ca3710f8929', 1698231003, ''),
(4, '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '993a5f9576ab63abfffbf45033ec9a0314777f0223cc4231d61ea1531045d2362472c82fe12de3c3a82e994f543568538b0875bae55d7ab379eed9203c0d898d', 1701073172, '');

-- --------------------------------------------------------

--
-- Table structure for table `accounts_meta`
--

CREATE TABLE `accounts_meta` (
  `id` int(11) NOT NULL,
  `accountNumber` char(36) NOT NULL,
  `bankNumber` char(36) NOT NULL,
  `accountHolderGivenName` text NOT NULL,
  `accountHolderFamilyName` text NOT NULL,
  `accountHolderDateOfBirth` text NOT NULL,
  `accountHolderIdentificationNumber` text NOT NULL,
  ` accountIdenificationType` text NOT NULL,
  `country` text NOT NULL,
  `accountHolderContactNumber1` text NOT NULL,
  `accountHolderContactNumber2` text DEFAULT NULL,
  `accountHolderEmailAddress` text NOT NULL,
  `accountHolderAddressLine1` text NOT NULL,
  `accountHolderAddressLine2` text DEFAULT NULL,
  `accountHolderAddressLine3` text DEFAULT NULL,
  `accountHolderPostalCode` text NOT NULL,
  `timestamp` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `accounts_meta`
--

INSERT INTO `accounts_meta` (`id`, `accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, ` accountIdenificationType`, `country`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode`, `timestamp`) VALUES
(1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'John', 'Doe', '1990-01-15', '123456789', '', '', '555-555-5555', '444-444-4444', 'johndoe@example.com', '123 Main St', 'Apt 4B', 'Building XYZ', '12345', 1698230686),
(2, 'a6392391-6d0e-4156-9a5c-d756e4c9347e', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'segun', 'martins', '1990-01-15', '123456555', '', '', '555-555-5555', '444-444-4444', 'johndoe@example.com', '123 Main St', 'Apt 4B', 'Building XYZ', '12345', 1698789465),
(3, '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'dexter', 'fizere', '1990-01-15', '1234565432112', '', '', '123456543', '1232123321', 'akanbiadebugba699@gmail.com', '34b punch street', '34b punch street', '34b punch street', '110221', 1701041936);

-- --------------------------------------------------------

--
-- Table structure for table `accounts_meta_usa`
--

CREATE TABLE `accounts_meta_usa` (
  `id` int(11) NOT NULL,
  `accounts_meta_id` int(11) NOT NULL,
  `us_bank_number` text NOT NULL,
  `sort_code` text NOT NULL,
  `creation_date` date NOT NULL DEFAULT current_timestamp(),
  `accounts_id` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `bank_account`
--

CREATE TABLE `bank_account` (
  `id` int(11) NOT NULL,
  `balance` float NOT NULL,
  `timestamp` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `bank_account`
--

INSERT INTO `bank_account` (`id`, `balance`, `timestamp`) VALUES
(1, 2.37, 1701073233);

-- --------------------------------------------------------

--
-- Table structure for table `beneficiaries`
--

CREATE TABLE `beneficiaries` (
  `beneficiaryId` int(11) NOT NULL,
  `userId` int(11) DEFAULT NULL,
  `fullName` varchar(100) DEFAULT NULL,
  `bankName` varchar(100) DEFAULT NULL,
  `bankAccountNumber` varchar(20) DEFAULT NULL,
  `bankRoutingNumber` varchar(20) DEFAULT NULL,
  `swiftCode` varchar(20) DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `biodata`
--

CREATE TABLE `biodata` (
  `id` int(11) NOT NULL,
  `surname` varchar(50) NOT NULL,
  `firstName` varchar(50) NOT NULL,
  `homeAddress` varchar(100) NOT NULL,
  `city` varchar(50) NOT NULL,
  `phoneNumber` varchar(15) NOT NULL,
  `dateOfBirth` date NOT NULL,
  `country` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `biodata`
--

INSERT INTO `biodata` (`id`, `surname`, `firstName`, `homeAddress`, `city`, `phoneNumber`, `dateOfBirth`, `country`) VALUES
(6, 'doe', 'John', '123 Main Street', 'New York', '+234 8088974888', '1990-01-15', 'USA'),
(7, 'adenugba', 'adeoluwa', '123 Main Street', 'lagos', '+234 8088974888', '1990-01-15', 'NGN'),
(11, 'adenugba', 'adeoluwa', '34b punch street arepo ogun state', 'ogun state', '+234 8088974888', '1990-01-15', 'NGN'),
(12, 'adenugba', 'adeoluwa', '34b punch street arepo ogun state', 'ogun state', '+234 8088974888', '1990-01-15', 'NGN'),
(13, 'adenugba', 'adeoluwa', '34b punch street arepo ogun state', 'ogun state', '+234 8088974888', '1990-01-15', 'NGN'),
(14, 'david', 'benjamin', '34b punch street arepo ogun state', 'ogun state', '+234 8088974998', '1990-01-15', 'NGN');

-- --------------------------------------------------------

--
-- Table structure for table `compliance_logs`
--

CREATE TABLE `compliance_logs` (
  `logId` int(11) NOT NULL,
  `userId` int(11) DEFAULT NULL,
  `action` varchar(100) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `details` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `identity`
--

CREATE TABLE `identity` (
  `id` int(11) NOT NULL,
  `bvn` varchar(11) DEFAULT NULL,
  `passport` varchar(50) DEFAULT NULL,
  `utilityBill` varchar(100) DEFAULT NULL,
  `picture` blob DEFAULT NULL,
  `userId` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `identity`
--

INSERT INTO `identity` (`id`, `bvn`, `passport`, `utilityBill`, `picture`, `userId`) VALUES
(4, '123456789', 'AB123456', 'UtilityBill123', 0x70726f66696c655f706963747572652e6a7067, 6),
(5, '123456789', 'AB123456', 'UtilityBill123', 0x70726f66696c655f706963747572652e6a7067, 7),
(9, '123456789', 'AB123456', 'UtilityBill123', 0x70726f66696c655f706963747572652e6a7067, 12),
(10, '123456789', 'AB123456', 'UtilityBill123', 0x70726f66696c655f706963747572652e6a7067, 13),
(11, '123456789', 'AB123456', 'UtilityBill123', 0x70726f66696c655f706963747572652e6a7067, 14);

-- --------------------------------------------------------

--
-- Table structure for table `kyc_documents`
--

CREATE TABLE `kyc_documents` (
  `documentId` int(11) NOT NULL,
  `userId` int(11) DEFAULT NULL,
  `documentType` enum('passport','national ID') DEFAULT NULL,
  `documentNumber` varchar(20) DEFAULT NULL,
  `documentImagePath` varchar(255) DEFAULT NULL,
  `expiryDate` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `privileges`
--

CREATE TABLE `privileges` (
  `privilege_id` int(11) NOT NULL,
  `role` varchar(50) NOT NULL,
  `privilege_name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `role_privileges`
--

CREATE TABLE `role_privileges` (
  `role_id` int(11) NOT NULL,
  `privilege_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `transactions`
--

CREATE TABLE `transactions` (
  `id` int(11) NOT NULL,
  `transaction` varchar(4) NOT NULL,
  `type` int(11) NOT NULL,
  `senderAccountNumber` char(36) NOT NULL,
  `senderBankNumber` char(36) NOT NULL,
  `receiverAccountNumber` char(36) NOT NULL,
  `receiverBankNumber` char(36) NOT NULL,
  `transactionAmount` float NOT NULL,
  `feeAmount` float NOT NULL,
  `narration` text NOT NULL,
  `timestamp` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `transactions`
--

INSERT INTO `transactions` (`id`, `transaction`, `type`, `senderAccountNumber`, `senderBankNumber`, `receiverAccountNumber`, `receiverBankNumber`, `transactionAmount`, `feeAmount`, `narration`, `timestamp`) VALUES
(177, 'pain', 1000, '123121212', '123121212', '01928918', '01928918', 1000, 0.1, '', 1698225211),
(178, 'pain', 1000, '123456', '123456', '123456', '123456', 1000, 0.1, '', 1698225939),
(179, 'pain', 1000, '123456', '123456', '123456', '123456', 1000, 0.1, '', 1698228222),
(180, 'pain', 1000, '122212', '122212', '123456', '123456', 1000, 0.1, '', 1698228292),
(181, 'pain', 1000, '123456', '123456', '122212', '122212', 1000, 0.1, '', 1698228364),
(182, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'befcedbd-0f53-48f4-8219-60477bffb9d6', 1000, 0.1, '', 1698231165),
(183, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '', 1000, 0.1, '', 1698233650),
(184, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '', 1000, 0.1, '', 1698233737),
(185, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '', 1000, 0.1, '', 1698233893),
(186, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', '', '', 1000, 0.1, '', 1698233985),
(187, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 1000, 0.1, '', 1698234144),
(188, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 200, 0.02, '', 1698236370),
(189, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 200, 0.02, '', 1698236616),
(190, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', 1000, 0.1, '', 1698789886),
(191, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', 1000, 0.1, '', 1698790122),
(192, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 1000, 0.1, '', 1698790224),
(193, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 300, 0.03, '', 1698790290),
(194, 'pain', 1000, 'befcedbd-0f53-48f4-8219-60477bffb9d6', 'a0299975-b8e2-4358-8f1a-911ee12dbaac', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', '', 300, 0.03, '', 1698790568),
(195, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', '', 200, 0.02, '', 1698790781),
(196, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', '', 200, 0.02, '', 1698953354),
(197, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', '', 200, 0.02, '', 1698953396),
(198, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', '', 200, 0.02, '', 1698953888),
(199, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', '', 100, 0.01, '', 1698953910),
(200, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', '', 100, 0.01, '', 1698953927),
(201, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 'a6392391-6d0e-4156-9a5c-d756e4c9347e', '', 100, 0.01, '', 1701041261),
(202, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 2500, 0.25, '', 1701042014),
(203, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 420, 0.042, '', 1701069320),
(204, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 420, 0.042, 'CR', 1701072239),
(205, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 420, 0.042, 'CR', 1701072303),
(206, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 420, 0.042, 'CR', 1701072835),
(207, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 420, 0.042, 'CR', 1701073021),
(208, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 1000, 0.1, 'CR', 1701073039),
(209, 'pain', 1, 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 1000, 0.1, 'CR', 1701073051),
(210, 'pain', 1, '2d192801-a08e-4cb8-bd86-0a8506a4fe6e', '', 'befcedbd-0f53-48f4-8219-60477bffb9d6', '', 1000, 0.1, 'DR', 1701073233);

-- --------------------------------------------------------

--
-- Table structure for table `v1accounts`
--

CREATE TABLE `v1accounts` (
  `id` int(11) NOT NULL,
  `userId` int(11) DEFAULT NULL,
  `accountHolderName` text NOT NULL,
  `type` enum('internal','uk','euro') DEFAULT NULL,
  `currencyCode` varchar(3) DEFAULT NULL,
  `balance` decimal(10,2) DEFAULT NULL,
  `sortCode` varchar(200) NOT NULL,
  `swiftCode` varchar(200) NOT NULL,
  `iban` varchar(200) NOT NULL,
  `overdraft` float NOT NULL,
  `availableBalance` float NOT NULL,
  `routingNumber` varchar(200) NOT NULL,
  `other` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `v1accounts`
--

INSERT INTO `v1accounts` (`id`, `userId`, `accountHolderName`, `type`, `currencyCode`, `balance`, `sortCode`, `swiftCode`, `iban`, `overdraft`, `availableBalance`, `routingNumber`, `other`) VALUES
(1, 13, '', 'internal', 'NGN', '0.00', '', '', '', 0, 0, '', NULL),
(2, 14, '', 'internal', 'NGN', '0.00', '', '', '', 0, 0, '', NULL);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `accounts`
--
ALTER TABLE `accounts`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `accounts_auth`
--
ALTER TABLE `accounts_auth`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `accounts_meta`
--
ALTER TABLE `accounts_meta`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `accounts_meta_usa`
--
ALTER TABLE `accounts_meta_usa`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `bank_account`
--
ALTER TABLE `bank_account`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `beneficiaries`
--
ALTER TABLE `beneficiaries`
  ADD PRIMARY KEY (`beneficiaryId`),
  ADD KEY `userId` (`userId`);

--
-- Indexes for table `biodata`
--
ALTER TABLE `biodata`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `compliance_logs`
--
ALTER TABLE `compliance_logs`
  ADD PRIMARY KEY (`logId`),
  ADD KEY `user_id` (`userId`);

--
-- Indexes for table `identity`
--
ALTER TABLE `identity`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`userId`);

--
-- Indexes for table `kyc_documents`
--
ALTER TABLE `kyc_documents`
  ADD PRIMARY KEY (`documentId`),
  ADD KEY `FK_userId` (`userId`);

--
-- Indexes for table `privileges`
--
ALTER TABLE `privileges`
  ADD PRIMARY KEY (`privilege_id`);

--
-- Indexes for table `role_privileges`
--
ALTER TABLE `role_privileges`
  ADD PRIMARY KEY (`role_id`,`privilege_id`),
  ADD KEY `privilege_id` (`privilege_id`);

--
-- Indexes for table `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `v1accounts`
--
ALTER TABLE `v1accounts`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `accounts`
--
ALTER TABLE `accounts`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `accounts_auth`
--
ALTER TABLE `accounts_auth`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `accounts_meta`
--
ALTER TABLE `accounts_meta`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `accounts_meta_usa`
--
ALTER TABLE `accounts_meta_usa`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `bank_account`
--
ALTER TABLE `bank_account`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `beneficiaries`
--
ALTER TABLE `beneficiaries`
  MODIFY `beneficiaryId` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `biodata`
--
ALTER TABLE `biodata`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- AUTO_INCREMENT for table `identity`
--
ALTER TABLE `identity`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- AUTO_INCREMENT for table `kyc_documents`
--
ALTER TABLE `kyc_documents`
  MODIFY `documentId` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `privileges`
--
ALTER TABLE `privileges`
  MODIFY `privilege_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `transactions`
--
ALTER TABLE `transactions`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=211;

--
-- AUTO_INCREMENT for table `v1accounts`
--
ALTER TABLE `v1accounts`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `beneficiaries`
--
ALTER TABLE `beneficiaries`
  ADD CONSTRAINT `beneficiaries_ibfk_1` FOREIGN KEY (`userId`) REFERENCES `biodata` (`id`);

--
-- Constraints for table `compliance_logs`
--
ALTER TABLE `compliance_logs`
  ADD CONSTRAINT `compliance_logs_ibfk_1` FOREIGN KEY (`userId`) REFERENCES `biodata` (`id`);

--
-- Constraints for table `identity`
--
ALTER TABLE `identity`
  ADD CONSTRAINT `identity_ibfk_1` FOREIGN KEY (`userId`) REFERENCES `biodata` (`id`);

--
-- Constraints for table `kyc_documents`
--
ALTER TABLE `kyc_documents`
  ADD CONSTRAINT `FK_userId` FOREIGN KEY (`userId`) REFERENCES `biodata` (`id`);

--
-- Constraints for table `role_privileges`
--
ALTER TABLE `role_privileges`
  ADD CONSTRAINT `role_privileges_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `accounts_auth` (`id`),
  ADD CONSTRAINT `role_privileges_ibfk_2` FOREIGN KEY (`privilege_id`) REFERENCES `privileges` (`privilege_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
