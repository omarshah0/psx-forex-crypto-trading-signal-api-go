-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "public";

-- CreateTable
CREATE TABLE "admins" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "admins_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "oauth_providers" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER NOT NULL,
    "provider" VARCHAR(50) NOT NULL,
    "provider_user_id" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "oauth_providers_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "schema_migrations" (
    "version" BIGINT NOT NULL,
    "dirty" BOOLEAN NOT NULL,

    CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version")
);

-- CreateTable
CREATE TABLE "trading_signals" (
    "id" SERIAL NOT NULL,
    "symbol" VARCHAR(50) NOT NULL,
    "stop_loss_price" DECIMAL(20,8) NOT NULL,
    "entry_price" DECIMAL(20,8) NOT NULL,
    "take_profit_price" DECIMAL(20,8) NOT NULL,
    "type" VARCHAR(10) NOT NULL,
    "result" VARCHAR(20),
    "return" DECIMAL(10,2),
    "created_by" INTEGER NOT NULL,
    "created_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "asset_class" VARCHAR(20),
    "duration_type" VARCHAR(20),
    "free_for_all" BOOLEAN DEFAULT false,
    "comments" TEXT,

    CONSTRAINT "trading_signals_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "users" (
    "id" SERIAL NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "blocked" BOOLEAN DEFAULT false,
    "created_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "hashed_password" VARCHAR(255),
    "email_verified" BOOLEAN DEFAULT false,
    "verification_token" VARCHAR(255),
    "verification_token_expires" TIMESTAMP(6),
    "reset_token" VARCHAR(255),
    "reset_token_expires" TIMESTAMP(6),
    "profile_picture" VARCHAR(1000),

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "packages" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "asset_class" VARCHAR(20) NOT NULL,
    "duration_type" VARCHAR(20) NOT NULL,
    "billing_cycle" VARCHAR(20) NOT NULL,
    "duration_days" INTEGER NOT NULL,
    "price" DECIMAL(10,2) NOT NULL,
    "description" TEXT,
    "is_active" BOOLEAN DEFAULT true,
    "created_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "packages_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "payment_history" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER NOT NULL,
    "package_id" INTEGER NOT NULL,
    "amount" DECIMAL(10,2) NOT NULL,
    "payment_method" VARCHAR(50),
    "payment_status" VARCHAR(20) NOT NULL,
    "transaction_id" VARCHAR(255),
    "metadata" JSONB,
    "created_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "payment_history_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "user_subscriptions" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER NOT NULL,
    "package_id" INTEGER NOT NULL,
    "price_paid" DECIMAL(10,2) NOT NULL,
    "subscribed_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "expires_at" TIMESTAMP(6) NOT NULL,
    "is_active" BOOLEAN DEFAULT true,
    "created_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "user_subscriptions_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "admins_user_id_key" ON "admins"("user_id");

-- CreateIndex
CREATE INDEX "idx_admins_user_id" ON "admins"("user_id");

-- CreateIndex
CREATE INDEX "idx_oauth_providers_email" ON "oauth_providers"("email");

-- CreateIndex
CREATE INDEX "idx_oauth_providers_provider" ON "oauth_providers"("provider");

-- CreateIndex
CREATE INDEX "idx_oauth_providers_user_id" ON "oauth_providers"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "oauth_providers_provider_provider_user_id_key" ON "oauth_providers"("provider", "provider_user_id");

-- CreateIndex
CREATE INDEX "idx_trading_signals_created_at" ON "trading_signals"("created_at" DESC);

-- CreateIndex
CREATE INDEX "idx_trading_signals_created_by" ON "trading_signals"("created_by");

-- CreateIndex
CREATE INDEX "idx_trading_signals_result" ON "trading_signals"("result");

-- CreateIndex
CREATE INDEX "idx_trading_signals_symbol" ON "trading_signals"("symbol");

-- CreateIndex
CREATE INDEX "idx_trading_signals_type" ON "trading_signals"("type");

-- CreateIndex
CREATE INDEX "idx_trading_signals_asset_class" ON "trading_signals"("asset_class");

-- CreateIndex
CREATE INDEX "idx_trading_signals_duration_type" ON "trading_signals"("duration_type");

-- CreateIndex
CREATE INDEX "idx_trading_signals_free_for_all" ON "trading_signals"("free_for_all");

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

-- CreateIndex
CREATE INDEX "idx_users_blocked" ON "users"("blocked");

-- CreateIndex
CREATE INDEX "idx_users_email" ON "users"("email");

-- CreateIndex
CREATE INDEX "idx_users_email_verified" ON "users"("email_verified");

-- CreateIndex
CREATE INDEX "idx_users_reset_token" ON "users"("reset_token");

-- CreateIndex
CREATE INDEX "idx_users_verification_token" ON "users"("verification_token");

-- CreateIndex
CREATE INDEX "idx_packages_asset_class" ON "packages"("asset_class");

-- CreateIndex
CREATE INDEX "idx_packages_billing_cycle" ON "packages"("billing_cycle");

-- CreateIndex
CREATE INDEX "idx_packages_duration_type" ON "packages"("duration_type");

-- CreateIndex
CREATE INDEX "idx_packages_is_active" ON "packages"("is_active");

-- CreateIndex
CREATE UNIQUE INDEX "packages_asset_class_duration_type_billing_cycle_key" ON "packages"("asset_class", "duration_type", "billing_cycle");

-- CreateIndex
CREATE INDEX "idx_payment_history_created_at" ON "payment_history"("created_at" DESC);

-- CreateIndex
CREATE INDEX "idx_payment_history_package_id" ON "payment_history"("package_id");

-- CreateIndex
CREATE INDEX "idx_payment_history_payment_status" ON "payment_history"("payment_status");

-- CreateIndex
CREATE INDEX "idx_payment_history_transaction_id" ON "payment_history"("transaction_id");

-- CreateIndex
CREATE INDEX "idx_payment_history_user_id" ON "payment_history"("user_id");

-- CreateIndex
CREATE INDEX "idx_user_subscriptions_expires_at" ON "user_subscriptions"("expires_at");

-- CreateIndex
CREATE INDEX "idx_user_subscriptions_is_active" ON "user_subscriptions"("is_active");

-- CreateIndex
CREATE INDEX "idx_user_subscriptions_package_id" ON "user_subscriptions"("package_id");

-- CreateIndex
CREATE INDEX "idx_user_subscriptions_user_active" ON "user_subscriptions"("user_id", "is_active", "expires_at");

-- CreateIndex
CREATE INDEX "idx_user_subscriptions_user_id" ON "user_subscriptions"("user_id");

-- AddForeignKey
ALTER TABLE "admins" ADD CONSTRAINT "admins_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "oauth_providers" ADD CONSTRAINT "oauth_providers_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "trading_signals" ADD CONSTRAINT "trading_signals_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "payment_history" ADD CONSTRAINT "payment_history_package_id_fkey" FOREIGN KEY ("package_id") REFERENCES "packages"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "payment_history" ADD CONSTRAINT "payment_history_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "user_subscriptions" ADD CONSTRAINT "user_subscriptions_package_id_fkey" FOREIGN KEY ("package_id") REFERENCES "packages"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "user_subscriptions" ADD CONSTRAINT "user_subscriptions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE NO ACTION;

