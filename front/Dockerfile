# Python 3.12.6 ベースイメージ
FROM python:3.12.6

# 環境変数設定
ENV PYTHONUNBUFFERED=1 \
    POETRY_HOME="/opt/poetry" \
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    POETRY_NO_INTERACTION=1 \
    PYTHONPATH="/"

# 必要なパッケージをインストール
RUN apt-get update && apt-get install -y curl && \
    curl -sSL https://install.python-poetry.org | python3 - && \
    ln -s $POETRY_HOME/bin/poetry /usr/local/bin/poetry

# 作業ディレクトリを front に設定
WORKDIR /front

# `pyproject.toml` と `poetry.lock` をコピー
COPY pyproject.toml poetry.lock ./

# 依存関係をインストール
RUN poetry install --no-root --only main

# ソースコードをコピー
COPY . .

# Streamlit のエントリーポイント
CMD ["poetry", "run", "streamlit", "run", "index.py", "--server.port=8501", "--server.address=0.0.0.0"]
