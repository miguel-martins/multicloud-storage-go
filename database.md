+------------------------+      +------------------------+      +------------------------+
|        Users           |      |        Files           |      |        Chunks          |
+------------------------+      +------------------------+      +------------------------+
| UserID (PK)            |      | FileID (PK)            |      | ChunkID (PK)           |
| Username               |      | UserID (FK)            |      | FileID (FK)            |
| PasswordHash           |      | FileName               |      | ChunkIndex             |
| PublicKey              |      | UploadDate             |      | AttributeEncryptedKey  |
| Additional Attributes  |      +------------------------+      | DeduplicationHash      |
+------------------------+                                      +------------------------+

+------------------------+      +------------------------+      +------------------------+
|  Attributes            |      | FileAttributes (Junction)|    |  CloudStorageProviders |
+------------------------+      +------------------------+      +------------------------+
| AttributeID (PK)       |      | FileID (FK)            |      | CloudStorageProviderID |
| AttributeName          |      | AttributeID (FK)       |      | ProviderName           |
| AttributeValue         |      +------------------------+      | EndpointURL            |
+------------------------+                                      | AccessKey              |
                                                                | SecretKey              |
                                                                +------------------------+

+------------------------+      +------------------------+
| ChunkCloudProviders    |      |       FileShares       |
+------------------------+      +------------------------+
| ChunkID (FK)           |      | ShareID (PK)           |
| CloudStorageProviderID |      | FileID (FK)            |
| CloudStorageLocation   |      | OwnerUserID (FK)       |
+------------------------+      | SharedWithUserID (FK)  |
                                +------------------------+
