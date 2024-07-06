#!/usr/bin/env bash

# Use este script para executar testes locais

RESULTS_WORKSPACE="$(pwd)/load-test/user-files/results"
GATLING_WORKSPACE="$(pwd)/load-test/user-files"
PROJECT_DIR="/root/gatling-maven-plugin-demo-scala"

runGatling() {
  cd $PROJECT_DIR || exit
  $GATLING_HOME gatling:test -Dgatling.simulationClass=RinhaBackendCrebitosSimulation \
        -Dgatling.simulation.description="Rinha de Backend - 2024/Q1: Crébito" \
        -Dgatling.resultsFolder=$RESULTS_WORKSPACE
}

startTest() {
    for i in {1..20}; do
        # 2 requests to wake the 2 api instances up :)
        curl --fail http://localhost:9999/clientes/1/extrato && \
        echo "" && \
        curl --fail http://localhost:9999/clientes/1/extrato && \
        echo "" && \
        runGatling && \
        break || sleep 2;
    done
}

startTest