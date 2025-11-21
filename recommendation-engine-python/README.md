# Motor de Recomendação de Conteúdos - Python FastAPI

Motor de recomendação integrado com o backend-go para gerar recomendações personalizadas baseadas em interações reais dos usuários.

## 🚀 Funcionalidades

- **Recomendações por Similaridade**: Usa Collaborative Filtering baseado em similaridade de cosseno entre usuários
- **Recomendações por Popularidade**: Baseado em ratings médios e número de interações
- **Integração com Banco de Dados**: Conecta ao mesmo banco do backend-go para usar interações reais
- **Atualização Incremental**: Modelo se atualiza automaticamente quando há novas interações
- **Fallback para Dados Simulados**: Usa dados simulados se o banco não estiver disponível

## 📋 Arquitetura

```
recommendation-engine-python/
├── app/
│   ├── api/
│   │   └── routes/
│   │       └── recommendations.py  # Endpoints da API
│   ├── core/
│   │   └── config.py               # Configurações
│   ├── models/
│   │   └── recommendation.py       # Modelo de recomendação
│   ├── schemas/
│   │   └── recommendation.py       # Schemas Pydantic
│   ├── services/
│   │   ├── database_service.py     # Conexão com banco de dados
│   │   └── dataset_service.py      # Gerenciamento de datasets
│   └── utils/
│       └── similarity.py           # Funções de similaridade
└── main.py                         # App FastAPI
```

## 🔧 Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` ou configure as seguintes variáveis:

```env
# Modo de dados: 'real' para usar banco de dados, 'simulated' para dados simulados
DATA_MODE=real

# Configurações do banco de dados (usado quando DATA_MODE=real)
DB_HOST=db
DB_PORT=3306
DB_USER=root
DB_PASSWORD=vertrigo
DB_NAME=content-recommender

# Ou use DATABASE_URL completa
# DATABASE_URL=mysql+pymysql://root:vertrigo@db:3306/content-recommender
```

## 📡 Endpoints da API

### 1. Obter Recomendações

```http
POST /recommendations/
Content-Type: application/json

{
  "user_id": 1,
  "top_n": 10,
  "method": "similarity"  // ou "popularity"
}
```

**Resposta:**
```json
{
  "user_id": 1,
  "recommendations": [
    {
      "content_id": 5,
      "score": 0.95,
      "title": "Conteúdo Recomendado"
    }
  ],
  "method": "similarity"
}
```

### 2. Criar Interação

```http
POST /recommendations/interactions
Content-Type: application/json

{
  "user_id": 1,
  "content_id": 5,
  "interaction_type": "like",  // "view", "like", "dislike", "rating"
  "rating": 5.0  // Opcional, usado quando interaction_type="rating"
}
```

**Resposta:**
```json
{
  "success": true,
  "message": "Interação criada com sucesso",
  "user_id": 1,
  "content_id": 5,
  "interaction_type": "like"
}
```

### 3. Estatísticas

```http
GET /recommendations/stats
```

**Resposta:**
```json
{
  "total_users": 150,
  "total_contents": 50,
  "total_interactions": 1250,
  "avg_interactions_per_user": 8.33,
  "data_mode": "real",
  "using_real_data": true
}
```

## 🔌 Integração com Backend-Go

O motor Python está integrado com o backend-go compartilhando o mesmo banco de dados MySQL. Quando o backend-go cria uma interação na tabela `user_interactions`, o motor Python pode:

1. **Buscar interações reais** do banco para treinar o modelo
2. **Receber novas interações** via API endpoint `/recommendations/interactions`
3. **Atualizar o modelo automaticamente** em background quando há novas interações

### Fluxo de Integração

```
App Mobile/Web
    ↓
Backend-Go (salva interação no banco)
    ↓
Motor Python (lê interações do banco OU recebe via API)
    ↓
Gera recomendações personalizadas
```

## 🐳 Executando com Docker

O serviço está configurado no `docker-compose.yml`:

```bash
docker-compose up recommender
```

O serviço estará disponível em `http://localhost:8000`

## 📚 Documentação Interativa

Acesse a documentação Swagger em:
- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc

## 🧪 Como Funciona

### 1. Modo Real (DATA_MODE=real)

- Conecta ao banco de dados MySQL compartilhado com o backend-go
- Busca interações da tabela `user_interactions`
- Busca conteúdos da tabela `contents`
- Modelo é treinado com dados reais
- Atualiza automaticamente quando há novas interações

### 2. Modo Simulado (DATA_MODE=simulated)

- Gera dados simulados para desenvolvimento/testes
- 100 usuários simulados, 50 conteúdos
- ~30% de densidade de interações
- Útil para testes sem banco de dados

## 🔄 Atualização do Modelo

O modelo é atualizado automaticamente quando:

1. **Nova interação via API**: Endpoint `/recommendations/interactions` atualiza o modelo em background
2. **Reinício do serviço**: Modelo é recarregado do banco de dados
3. **Atualização manual**: (próxima versão) endpoint para forçar recarregamento

## 📊 Métodos de Recomendação

### Similarity (Collaborative Filtering)
- Encontra usuários similares usando similaridade de cosseno
- Recomenda conteúdos que usuários similares gostaram
- Melhor para personalização individual

### Popularity
- Baseado em ratings médios e número de interações
- Recomenda conteúdos mais populares
- Melhor para usuários novos (cold start)

## 🛠️ Tecnologias

- **FastAPI**: Framework web moderno e rápido
- **Pandas**: Manipulação de dados
- **Scikit-Learn**: Cálculo de similaridade
- **SQLAlchemy**: ORM para banco de dados
- **PyMySQL**: Driver MySQL

## 📝 Notas

- O modelo recarrega automaticamente quando há novas interações (em background)
- Para grandes volumes, considere implementar cache ou recarregamento periódico
- O modelo atual é baseado em memória - para produção, considere persistir o modelo treinado

