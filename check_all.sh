echo "===== DOCKER COMPOSE PS ====="
docker compose ps

echo
echo "===== DOCKER LOGS (last 20 lines) ====="
containers=(api worker mongo cassandra temporal temporal-ui)
for c in "${containers[@]}"; do
  echo
  echo "----- $c -----"
  docker compose logs --tail=20 $c
done