"""
Serviço para conectar e buscar dados do banco de dados
Integra com o backend-go para acessar interações e conteúdos reais
"""
from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker, Session
from sqlalchemy.exc import SQLAlchemyError
import pandas as pd
from typing import Tuple, Optional, List, Dict
import logging
from app.core.config import settings

logger = logging.getLogger(__name__)

class DatabaseService:
    """Serviço para gerenciar conexão com banco de dados"""
    
    def __init__(self):
        self.engine = None
        self.SessionLocal = None
        self._initialize_connection()
    
    def _initialize_connection(self):
        """Inicializa a conexão com o banco de dados"""
        try:
            database_url = settings.get_database_url()
            self.engine = create_engine(
                database_url,
                pool_pre_ping=True,  # Verifica conexão antes de usar
                pool_recycle=3600,   # Recicla conexões a cada hora
                echo=False
            )
            self.SessionLocal = sessionmaker(
                autocommit=False,
                autoflush=False,
                bind=self.engine
            )
            logger.info(f"Conectado ao banco de dados: {settings.db_host}:{settings.db_port}/{settings.db_name}")
        except Exception as e:
            logger.error(f"Erro ao conectar ao banco de dados: {e}")
            self.engine = None
            self.SessionLocal = None
    
    def get_session(self) -> Optional[Session]:
        """Retorna uma sessão do banco de dados"""
        if self.SessionLocal is None:
            return None
        return self.SessionLocal()
    
    def is_connected(self) -> bool:
        """Verifica se está conectado ao banco de dados"""
        if self.engine is None:
            return False
        try:
            with self.engine.connect() as conn:
                conn.execute(text("SELECT 1"))
            return True
        except Exception as e:
            logger.error(f"Erro ao verificar conexão: {e}")
            return False
    
    def fetch_interactions(self) -> Optional[pd.DataFrame]:
        """
        Busca interações reais do banco de dados
        
        Returns:
            pd.DataFrame: DataFrame com colunas user_id, content_id, rating, interaction_type
        """
        if not self.is_connected():
            logger.warning("Não conectado ao banco de dados")
            return None
        
        try:
            query = text("""
                SELECT 
                    user_id,
                    content_id,
                    CASE 
                        WHEN interaction_type = 'like' THEN 5.0
                        WHEN interaction_type = 'dislike' THEN 1.0
                        WHEN interaction_type = 'view' THEN 3.0
                        WHEN rating IS NOT NULL THEN rating
                        ELSE 3.0
                    END as rating,
                    interaction_type
                FROM user_interactions
                ORDER BY created_at DESC
            """)
            
            df = pd.read_sql(query, self.engine)
            
            if df.empty:
                logger.info("Nenhuma interação encontrada no banco de dados")
                return None
            
            # Converter user_id e content_id para int
            df['user_id'] = df['user_id'].astype(int)
            df['content_id'] = df['content_id'].astype(int)
            df['rating'] = df['rating'].astype(float)
            
            logger.info(f"Carregadas {len(df)} interações do banco de dados")
            return df
            
        except SQLAlchemyError as e:
            logger.error(f"Erro ao buscar interações: {e}")
            return None
    
    def fetch_contents(self) -> Optional[pd.DataFrame]:
        """
        Busca conteúdos reais do banco de dados
        
        Returns:
            pd.DataFrame: DataFrame com colunas content_id, title, description
        """
        if not self.is_connected():
            logger.warning("Não conectado ao banco de dados")
            return None
        
        try:
            query = text("""
                SELECT 
                    id as content_id,
                    title,
                    description,
                    type as content_type
                FROM contents
                ORDER BY id
            """)
            
            df = pd.read_sql(query, self.engine)
            
            if df.empty:
                logger.info("Nenhum conteúdo encontrado no banco de dados")
                return None
            
            df['content_id'] = df['content_id'].astype(int)
            
            logger.info(f"Carregados {len(df)} conteúdos do banco de dados")
            return df
            
        except SQLAlchemyError as e:
            logger.error(f"Erro ao buscar conteúdos: {e}")
            return None
    
    def create_interaction(
        self,
        user_id: int,
        content_id: int,
        interaction_type: str,
        rating: Optional[float] = None
    ) -> bool:
        """
        Cria uma nova interação no banco de dados
        
        Args:
            user_id: ID do usuário
            content_id: ID do conteúdo
            interaction_type: Tipo de interação ('view', 'like', 'dislike', etc.)
            rating: Rating opcional (1-5)
            
        Returns:
            bool: True se criado com sucesso, False caso contrário
        """
        if not self.is_connected():
            logger.warning("Não conectado ao banco de dados")
            return False
        
        try:
            query = text("""
                INSERT INTO user_interactions 
                (user_id, content_id, interaction_type, rating, created_at)
                VALUES (:user_id, :content_id, :interaction_type, :rating, NOW())
            """)
            
            with self.engine.begin() as conn:
                conn.execute(query, {
                    'user_id': user_id,
                    'content_id': content_id,
                    'interaction_type': interaction_type,
                    'rating': rating
                })
            
            logger.info(f"Interação criada: user_id={user_id}, content_id={content_id}, type={interaction_type}")
            return True
            
        except SQLAlchemyError as e:
            logger.error(f"Erro ao criar interação: {e}")
            return False
    
    def get_user_interactions_count(self, user_id: int) -> int:
        """Retorna o número de interações de um usuário"""
        if not self.is_connected():
            return 0
        
        try:
            query = text("SELECT COUNT(*) as count FROM user_interactions WHERE user_id = :user_id")
            result = pd.read_sql(query, self.engine, params={'user_id': user_id})
            return int(result.iloc[0]['count']) if not result.empty else 0
        except Exception as e:
            logger.error(f"Erro ao contar interações: {e}")
            return 0

# Singleton
database_service = DatabaseService()

